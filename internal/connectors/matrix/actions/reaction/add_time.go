package reaction

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

// AddTimeAction adds time to an event based on reactions.
type AddTimeAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
}

// Configure is called on startup and sets all dependencies.
func (action *AddTimeAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
}

// Name of the action.
func (action *AddTimeAction) Name() string {
	return "Add Time"
}

// GetDocu returns the documentation for the action.
func (action *AddTimeAction) GetDocu() (title, explaination string, examples []string) {
	return "Add Time",
		`Use the following reactions to add some more hours to an event: 1️⃣, 2️⃣, 3️⃣, 4️⃣, 5️⃣, 6️⃣, 7️⃣, 8️⃣, 9️⃣, 🔟.
Or use ▶️/⏩ to move the event to tomorrow/next week.`,
		[]string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟", "➕"}
}

// Selector defines on which reactions this action should be called.
func (action *AddTimeAction) Selector() []string {
	return []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟", "➕", "⏩", "▶️"}
}

// HandleEvent is where the reaction event and the related message get's send to if it matches the Selector.
func (action *AddTimeAction) HandleEvent(event *matrix.ReactionEvent, reactionToMessage *matrixdb.MatrixMessage) {
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
	evt.Active = true
	action.addTimeToEvent(event.Content.RelatesTo.Key, evt)

	_, err := action.db.UpdateEvent(evt)
	if err != nil {
		l.Err(err)
		_ = action.messenger.SendMessageAsync(messenger.PlainTextMessage(
			"Whoopsie, can not update the event as requested.",
			event.Room.RoomID,
		))
		return
	}

	err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
		fmt.Sprintf(
			"I rescheduled this reminder to %s!",
			format.ToLocalTime(evt.Time, event.Room.TimeZone),
		),
		reactionToMessage.ID,
		reactionToMessage.Body,
		event.Event.Sender.String(),
		event.Room.RoomID,
	))
	if err != nil {
		l.Err(err)
	}
}

func (action *AddTimeAction) addTimeToEvent(reactionKey string, event *database.Event) {
	switch reactionKey {
	case "1️⃣":
		event.Time = time.Now().Add(1 * time.Hour)
	case "2️⃣":
		event.Time = time.Now().Add(2 * time.Hour)
	case "3️⃣":
		event.Time = time.Now().Add(3 * time.Hour)
	case "4️⃣":
		event.Time = time.Now().Add(4 * time.Hour)
	case "5️⃣":
		event.Time = time.Now().Add(5 * time.Hour)
	case "6️⃣":
		event.Time = time.Now().Add(6 * time.Hour)
	case "7️⃣":
		event.Time = time.Now().Add(7 * time.Hour)
	case "8️⃣":
		event.Time = time.Now().Add(8 * time.Hour)
	case "9️⃣":
		event.Time = time.Now().Add(9 * time.Hour)
	case "🔟":
		event.Time = time.Now().Add(10 * time.Hour)
	case "▶️":
		event.Time = event.Time.Add(24 * time.Hour)
		// Make sure event time is always in the future.
		for time.Until(event.Time) <= 0 {
			event.Time = event.Time.Add(24 * time.Hour)
		}
	case "⏩":
		event.Time = event.Time.Add(7 * 24 * time.Hour)
		// Make sure event time is always in the future.
		for time.Until(event.Time) <= 0 {
			event.Time = event.Time.Add(7 * 24 * time.Hour)
		}
	default:
		action.logger.Errorf("do not know what time to add for key '%s'", reactionKey)
	}
}
