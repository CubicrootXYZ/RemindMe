package reply

import (
	"regexp"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

var deleteEventActionRegex = regexp.MustCompile("(?i)^(delete|remove|cancel)[ ]*$")

// DeleteEventAction deletes an event.
type DeleteEventAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
}

// Configure is called on startup and sets all dependencies.
func (action *DeleteEventAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
}

// Name of the action
func (action *DeleteEventAction) Name() string {
	return "Delete Event"
}

// GetDocu returns the documentation for the action.
func (action *DeleteEventAction) GetDocu() (title, explaination string, examples []string) {
	return "Delete Event",
		"Delete an Event by replying to it",
		[]string{"delete", "remove", "cancel"}
}

// Selector defines a regex on what messages the action should be used.
func (action *DeleteEventAction) Selector() *regexp.Regexp {
	return deleteEventActionRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *DeleteEventAction) HandleEvent(event *matrix.MessageEvent, replyToMessage *matrixdb.MatrixMessage) {
	if replyToMessage.EventID == nil || replyToMessage.Event == nil {
		// No event given, can not update anything
		action.logger.Debugf("can not delete event with event ID nil")
		return
	}

	err := action.db.DeleteEvent(replyToMessage.Event)
	if err != nil {
		action.logger.Errorf("failed to update event in database: %w", err)
		return
	}

	_ = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
		"I deleted this Event!",
		event.Event.ID.String(),
		event.Content.Body,
		event.Event.Sender.String(),
		event.Room.RoomID,
	))

	// Best effort approach to delete messages related to that event.
	messages, err := action.matrixDB.ListMessages(matrixdb.ListMessageOpts{
		RoomID:  &event.Room.ID,
		EventID: replyToMessage.EventID,
	})
	if err != nil {
		action.logger.Errorf("failed to list messages for event: %w", err)
		return
	}

	for _, message := range messages {
		err := action.messenger.DeleteMessageAsync(&messenger.Delete{
			ExternalIdentifier:        message.ID,
			ChannelExternalIdentifier: event.Room.RoomID,
		})
		if err != nil {
			action.logger.Errorf("failed to delete message: %w", err)
		}
	}
}
