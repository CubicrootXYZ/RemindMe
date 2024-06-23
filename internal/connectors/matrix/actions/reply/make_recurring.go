package reply

import (
	"fmt"
	"regexp"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/gonaturalduration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mapping"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/msghelper"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

var makeRecurringActionRegex = regexp.MustCompile("(?i)^(repeat|every|each|always|recurring|all|any).*(second|minute|hour|day|week|month|year)(|s)[ ]*$")

// MakeRecurringAction makes an event recurring.
type MakeRecurringAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
	storer    *msghelper.Storer
}

// Configure is called on startup and sets all dependencies.
func (action *MakeRecurringAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
	action.storer = msghelper.NewStorer(matrixDB, messenger, logger)
}

// Name of the action
func (action *MakeRecurringAction) Name() string {
	return "Make Event Recurring"
}

// GetDocu returns the documentation for the action.
func (action *MakeRecurringAction) GetDocu() (title, explaination string, examples []string) {
	return "Make Event Recurring",
		"Make an event recurring by replying with a duration.",
		[]string{"every 10 days", "each twenty two hours and five seconds"}
}

// Selector defines a regex on what messages the action should be used.
func (action *MakeRecurringAction) Selector() *regexp.Regexp {
	return makeRecurringActionRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *MakeRecurringAction) HandleEvent(event *matrix.MessageEvent, replyToMessage *matrixdb.MatrixMessage) {
	if replyToMessage.EventID == nil || replyToMessage.Event == nil {
		// No event given, can not update anything
		action.logger.Debugf("can not change event with event ID nil")
		return
	}

	// Get duration from message
	duration := gonaturalduration.ParseNumber(event.Content.Body)
	if duration <= time.Minute {
		action.logger.Infof("missing duration in message")
		_ = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Sorry I was not able to understand the duration from this message.",
			event.Event.ID.String(),
			event.Content.Body,
			event.Event.Sender.String(),
			event.Room.RoomID,
		))
		return
	}

	message := mapping.MessageFromEvent(event)
	message.Type = matrixdb.MessageTypeChangeEvent
	message.EventID = replyToMessage.EventID
	_, err := action.matrixDB.NewMessage(message)
	if err != nil {
		action.logger.Errorf("failed to save message to db: %v", err)
		return
	}

	defaultRepeatUntil := time.Now().Add((5 * 365 * 24 * time.Hour))
	replyToMessage.Event.RepeatInterval = &duration
	replyToMessage.Event.RepeatUntil = &defaultRepeatUntil

	dbEvent, err := action.db.UpdateEvent(replyToMessage.Event)
	if err != nil {
		action.logger.Err(err)
		_ = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Sorry, that failed :/.",
			event.Event.ID.String(),
			event.Content.Body,
			event.Event.Sender.String(),
			event.Room.RoomID,
		))
		return
	}

	msg := fmt.Sprintf(
		"Updated the event to remind you every %s until %s",
		format.ToNiceDuration(duration),
		format.ToLocalTime(*dbEvent.RepeatUntil, event.Room.TimeZone),
	)
	go action.storer.SendAndStoreResponse(
		msg,
		matrixdb.MessageTypeTimezoneChange,
		*event,
		msghelper.WithEventID(*replyToMessage.EventID),
	)
}
