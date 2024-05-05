package message

import (
	"regexp"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/msghelper"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"

	_ "time/tzdata" // Import timezone data.
)

var changeTimezoneActionRegex = regexp.MustCompile("(?i)^set timezone .*$")
var timezoneCaptureGroup = regexp.MustCompile("(?i)^set timezone (.*)$")

// ChangeTimezoneAction allows setting a timezone for the matrix channel.
type ChangeTimezoneAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
	storer    *msghelper.Storer
}

// Configure is called on startup and sets all dependencies.
func (action *ChangeTimezoneAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
	action.storer = msghelper.NewStorer(matrixDB, messenger, logger)
}

// Name of the action
func (action *ChangeTimezoneAction) Name() string {
	return "Change Timezone"
}

// GetDocu returns the documentation for the action.
func (action *ChangeTimezoneAction) GetDocu() (title, explaination string, examples []string) {
	return "Change Timezone",
		"Change the timezone for this channel",
		[]string{"set timezone Europe/Berlin", "set timezone America/New_York", "set timezone Asia/Shanghai"}
}

// Selector defines a regex on what messages the action should be used.
func (action *ChangeTimezoneAction) Selector() *regexp.Regexp {
	return changeTimezoneActionRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *ChangeTimezoneAction) HandleEvent(event *matrix.MessageEvent) {
	tz := ""
	matches := timezoneCaptureGroup.FindStringSubmatch(event.Content.Body)
	if len(matches) >= 2 {
		tz = matches[1]
	}
	_, err := time.LoadLocation(tz)
	if err != nil {
		action.logger.Infof("failed to load timezone '%s' with: %s", tz, err.Error())
		action.storer.SendAndStoreResponse("Sorry, but I do not know what timezone this is.", matrixdb.MessageTypeTimezoneChange, *event)
		return
	}

	room := event.Room
	tzBefore := room.TimeZone
	room.TimeZone = tz

	_, err = action.matrixDB.UpdateRoom(room)
	if err != nil {
		action.logger.Err(err)
		err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Ups, that did not work 😨",
			event.Event.ID.String(),
			event.Content.Body,
			event.Event.Sender.String(),
			event.Room.RoomID,
		))
		if err != nil {
			action.logger.Err(err)
		}
		return
	}

	msgBuilder := format.Formater{}
	msgBuilder.Text("Changed this channels timezone from ")
	if tzBefore == "" {
		tzBefore = "UTC"
	}
	msgBuilder.Italic(tzBefore)
	msgBuilder.Text(" to ")
	msgBuilder.Text(tz)
	msgBuilder.Text(" 🛫 🛬")

	msg, _ := msgBuilder.Build()
	go action.storer.SendAndStoreResponse(msg, matrixdb.MessageTypeTimezoneChange, *event)
}
