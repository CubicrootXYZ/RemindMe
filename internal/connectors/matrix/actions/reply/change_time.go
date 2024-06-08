package reply

import (
	"fmt"
	"regexp"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mapping"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/msghelper"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

var changeTimeRegex = regexp.MustCompile(".*")

// ChangeTimeAction acts as a template for new actions.
type ChangeTimeAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
	storer    *msghelper.Storer
}

// Configure is called on startup and sets all dependencies.
func (action *ChangeTimeAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
	action.storer = msghelper.NewStorer(matrixDB, messenger, logger)
}

// Name of the action
func (action *ChangeTimeAction) Name() string {
	return "Change time"
}

// GetDocu returns the documentation for the action.
func (action *ChangeTimeAction) GetDocu() (title, explaination string, examples []string) {
	return "Change time",
		"Change the time of a reminder by replying to a reminder message.",
		[]string{"January, 5th", "at 5 pm", "tomorrow"}
}

// Selector defines a regex on what messages the action should be used.
func (action *ChangeTimeAction) Selector() *regexp.Regexp {
	return changeTimeRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *ChangeTimeAction) HandleEvent(event *matrix.MessageEvent, replyToMessage *matrixdb.MatrixMessage) {
	if replyToMessage.EventID == nil || replyToMessage.Event == nil {
		// No event given, can not update anything
		action.logger.Debugf("can not update event with event ID nil")
		return
	}

	remindTime, err := format.ParseTime(event.Content.Body, event.Room.TimeZone, false)
	if err != nil {
		action.logger.Err(err)
		_ = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Sorry I was not able to understand the remind date and time from this message.",
			event.Event.ID.String(),
			event.Content.Body,
			event.Event.Sender.String(),
			event.Room.RoomID,
		))
		return
	}

	message := mapping.MessageFromEvent(event)
	message.Type = matrixdb.MessageTypeChangeEvent
	_, err = action.matrixDB.NewMessage(message)
	if err != nil {
		action.logger.Errorf("failed to save message to db: %v", err)
		return
	}

	replyToMessage.Event.Time = remindTime
	_, err = action.db.UpdateEvent(replyToMessage.Event)
	if err != nil {
		action.logger.Errorf("failed to update event in database: %v", err)
		return
	}

	go action.storer.SendAndStoreResponse(
		fmt.Sprintf("I rescheduled your reminder \"%s\" to %s.", replyToMessage.Event.Message, format.ToLocalTime(replyToMessage.Event.Time, event.Room.TimeZone)),
		matrixdb.MessageTypeChangeEvent,
		*event,
		msghelper.WithEventID(replyToMessage.Event.ID),
	)
}
