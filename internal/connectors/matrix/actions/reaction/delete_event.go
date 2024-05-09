package reaction

import (
	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

// DeleteEventAction takes cafe of delete requests via reactions.
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

// Name of the action.
func (action *DeleteEventAction) Name() string {
	return "Delete Event"
}

// GetDocu returns the documentation for the action.
func (action *DeleteEventAction) GetDocu() (title, explaination string, examples []string) {
	return "Delete Event",
		"React with a ❌ to delete an event.",
		[]string{"❌"}
}

// Selector defines on which reactions this action should be called.
func (action *DeleteEventAction) Selector() []string {
	return []string{"❌"}
}

// HandleEvent is where the reaction event and the related message get's send to if it matches the Selector.
func (action *DeleteEventAction) HandleEvent(event *matrix.ReactionEvent, reactionToMessage *matrixdb.MatrixMessage) {
	l := action.logger.WithFields(
		map[string]any{
			"reaction":        event.Content.RelatesTo.Key,
			"room":            reactionToMessage.RoomID,
			"related_message": reactionToMessage.ID,
			"user":            event.Event.Sender,
		},
	)
	if reactionToMessage.EventID == nil || reactionToMessage.Event == nil {
		l.Infof("skipping because message does not relate to any event")
		return
	}

	evt := reactionToMessage.Event
	evt.Active = false

	_, err := action.db.UpdateEvent(evt)
	if err != nil {
		l.Err(err)
		_ = action.messenger.SendMessageAsync(messenger.PlainTextMessage(
			"Whoopsie, can not delete the event as requested.",
			event.Room.RoomID,
		))
		return
	}

	// Later improvement: extend storer for reaction events and store this
	// message to the database as well.
	err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
		"Deleted event \""+reactionToMessage.Event.Message+"\"",
		reactionToMessage.ID,
		reactionToMessage.Body,
		event.Event.Sender.String(),
		event.Room.RoomID,
	))
	if err != nil {
		action.logger.Errorf("failed to send message: %w", err)
		return
	}

	// Best effort approach to delete messages related to that event.
	messages, err := action.matrixDB.ListMessages(matrixdb.ListMessageOpts{
		RoomID:  &event.Room.ID,
		EventID: reactionToMessage.EventID,
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
