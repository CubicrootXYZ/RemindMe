package message

import (
	"log/slog"
	"regexp"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/msghelper"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

var changeEventActionRegex = regexp.MustCompile("(?i)(^(change|update|set)[ ]+(reminder|reminder id|)[ ]*[0-9]+)")

// ChangeEventAction for changing an existing event.
type ChangeEventAction struct {
	logger    *slog.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
	storer    *msghelper.Storer
}

// Configure is called on startup and sets all dependencies.
func (action *ChangeEventAction) Configure(logger *slog.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
	action.storer = msghelper.NewStorer(matrixDB, messenger, logger)
}

// Name of the action
func (action *ChangeEventAction) Name() string {
	return "Change event"
}

func (action *ChangeEventAction) GetDocu() (title, explaination string, examples []string) {
	return "Change event",
		"Change an existing Event.",
		[]string{"change reminder 1 to tomorrow", "update 68 to Saturday 4 pm"}
}

// Selector defines a regex on what messages the action should be used.
func (action *ChangeEventAction) Selector() *regexp.Regexp {
	return changeEventActionRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *ChangeEventAction) HandleEvent(event *matrix.MessageEvent) {
	match := changeEventActionRegex.Find([]byte(event.Content.Body))
	if match == nil {
		go action.storer.SendAndStoreResponse(
			"Ups, seems like there is a reminder ID missing in your message.",
			matrixdb.MessageTypeChangeEventError,
			*event,
		)
		return
	}

	eventID, err := format.GetSuffixInt(string(match))
	if err != nil {
		go action.storer.SendAndStoreResponse(
			"Ups, seems like there is a reminder ID missing in your message.",
			matrixdb.MessageTypeChangeEventError,
			*event,
		)
		return
	}

	newTime, err := format.ParseTime(event.Content.Body, event.Room.TimeZone, false)
	if err != nil {
		action.logger.Error("failed to parse time", "error", err)
		go action.storer.SendAndStoreResponse(
			"Ehm, sorry to say that, but I was not able to understand the time to schedule the reminder to.",
			matrixdb.MessageTypeChangeEventError,
			*event,
		)
		return
	}

	events, err := action.db.ListEvents(&database.ListEventsOpts{
		IDs:       []uint{uint(eventID)},
		ChannelID: &event.Channel.ID,
	})
	if err != nil || len(events) == 0 {
		if err != nil {
			action.logger.Error("failed to list events", "error", err)
		}
		go action.storer.SendAndStoreResponse(
			"This reminder is not in my database.",
			matrixdb.MessageTypeChangeEventError,
			*event,
		)
		return
	}

	evt := &events[0]
	evt.Time = newTime
	evt, err = action.db.UpdateEvent(evt)
	if err != nil || len(events) == 0 {
		action.logger.Error("failed to update event", "error", err)
		go action.storer.SendAndStoreResponse(
			"Whups, this did not work, sorry.",
			matrixdb.MessageTypeChangeEventError,
			*event,
		)
		return
	}

	msgFormater := format.Formater{}
	msgFormater.TextLine("I rescheduled your reminder")
	msgFormater.QuoteLine(evt.Message)
	msgFormater.Text("to ")
	msgFormater.Text(format.ToLocalTime(newTime, event.Room.TimeZone))
	msg, formattedMsg := msgFormater.Build()

	go action.storer.SendAndStoreMessage(msg, formattedMsg, matrixdb.MessageTypeChangeEvent, *event)
}
