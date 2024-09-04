package reaction

import (
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

// RescheduleRepeatingAction takes care of rescheduling a repeating event.
type RescheduleRepeatingAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
}

// Configure is called on startup and sets all dependencies.
func (action *RescheduleRepeatingAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
}

// Name of the action.
func (action *RescheduleRepeatingAction) Name() string {
	return "Reschedule Repeating Event"
}

// GetDocu returns the documentation for the action.
func (action *RescheduleRepeatingAction) GetDocu() (title, explaination string, examples []string) {
	return "Reschedule Repeating Event",
		"React with a 🔂 to get reminded again in 1 hour without changing the next repeat cycle.",
		[]string{"🔂"}
}

// Selector defines on which reactions this action should be called.
func (action *RescheduleRepeatingAction) Selector() []string {
	return []string{"🔂"}
}

// HandleEvent is where the reaction event and the related message get's send to if it matches the Selector.
func (action *RescheduleRepeatingAction) HandleEvent(event *matrix.ReactionEvent, reactionToMessage *matrixdb.MatrixMessage) {
	// TODO tests & add new reaction
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

	// Clone the event without being repetitive.
	newEvt := &database.Event{
		Time:      time.Now().Add(time.Hour),
		Duration:  reactionToMessage.Event.Duration,
		Message:   reactionToMessage.Event.Message,
		Active:    true,
		ChannelID: reactionToMessage.Event.ChannelID,
		InputID:   reactionToMessage.Event.InputID,
	}
	_, err := action.db.NewEvent(newEvt)
	if err != nil {
		l.Err(err)
		_ = action.messenger.SendMessageAsync(messenger.PlainTextMessage(
			"Whoopsie, can not update the event as requested.",
			event.Room.RoomID,
		))
		return
	}

	err = action.messenger.DeleteMessageAsync(&messenger.Delete{
		ExternalIdentifier:        reactionToMessage.ID,
		ChannelExternalIdentifier: reactionToMessage.Room.RoomID,
	})
	if err != nil {
		l.Err(err)
	}
}
