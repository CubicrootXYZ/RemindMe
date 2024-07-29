package message

import (
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

var listEventsActionRegex = regexp.MustCompile("(?i)^((list|show)(| all| the)(| reminders| my reminders)(| please)|^reminders|^reminder)[ ]*$")

// ListEventsAction lists all events.
type ListEventsAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
	storer    *msghelper.Storer
}

// Configure is called on startup and sets all dependencies.
func (action *ListEventsAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
	action.storer = msghelper.NewStorer(matrixDB, messenger, logger)
}

// Name of the action
func (action *ListEventsAction) Name() string {
	return "List All Events"
}

// GetDocu returns the documentation for the action.
func (action *ListEventsAction) GetDocu() (title, explaination string, examples []string) {
	return "List All Events",
		"List all events in this channel.",
		[]string{"list", "list reminders", "show", "show reminders", "list my reminders", "reminders"}
}

// Selector defines a regex on what messages the action should be used.
func (action *ListEventsAction) Selector() *regexp.Regexp {
	return listEventsActionRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *ListEventsAction) HandleEvent(event *matrix.MessageEvent) {
	events, err := action.db.ListEvents(&database.ListEventsOpts{
		ChannelID: &event.Channel.ID,
	})
	if err != nil {
		err2 := action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"There was an issue accessing the data 🤨",
			event.Event.ID.String(),
			event.Content.Body,
			event.Event.Sender.String(),
			event.Room.RoomID,
		))
		if err2 != nil {
			action.logger.Err(err2)
		}
		action.logger.Err(err)
		return
	}

	message := mapping.MessageFromEvent(event)
	message.Type = matrixdb.MessageTypeEventList
	_, err = action.matrixDB.NewMessage(message)
	if err != nil {
		action.logger.Err(err)
	}

	msg := format.Formater{}
	msg.Title("Your Events")

	msgEvts, msgEvtsFormatted := format.InfoFromEvents(events, event.Room.TimeZone)

	msgOut, msgOutFormatted := msg.Build()

	go action.storer.SendAndStoreMessage(msgOut+msgEvts, msgOutFormatted+msgEvtsFormatted, matrixdb.MessageTypeEventList, *event)
}
