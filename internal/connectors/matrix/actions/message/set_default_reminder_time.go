package message

import (
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mapping"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/msghelper"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

var setDefaultReminderTimeRegex = regexp.MustCompile("(?i)^(set|update|change|)[ ]*(the|a|my|)[ ]*(default|)[ ]*(reminder time|remind time|reminder).*")

// SetDefaultReminderTimeAction sets the time for the daily reminder.
type SetDefaultReminderTimeAction struct {
	logger    *slog.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
	storer    *msghelper.Storer
}

// Configure is called on startup and sets all dependencies.
func (action *SetDefaultReminderTimeAction) Configure(logger *slog.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
	action.storer = msghelper.NewStorer(matrixDB, messenger, logger)
}

// Name of the action
func (action *SetDefaultReminderTimeAction) Name() string {
	return "Set Default Reminder Time"
}

// GetDocu returns the documentation for the action.
func (action *SetDefaultReminderTimeAction) GetDocu() (title, explanation string, examples []string) {
	return "Set Default Reminder Time",
		"Set the default time to use if no time is given",
		[]string{"default reminder at 9am", "default reminder at 13:00", "set the remind time to 4:00"}
}

// Selector defines a regex on what messages the action should be used.
func (action *SetDefaultReminderTimeAction) Selector() *regexp.Regexp {
	return setDefaultReminderTimeRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *SetDefaultReminderTimeAction) HandleEvent(event *matrix.MessageEvent) {
	dbMsg := mapping.MessageFromEvent(event)
	dbMsg.Type = matrixdb.MessageTypeSetDefaultReminderTime

	_, err := action.matrixDB.NewMessage(dbMsg)
	if err != nil {
		action.logger.Error("failed to store message to database", "error", err)
	}

	timeRemind, err := format.ParseTime(event.Channel, event.Content.Body, event.Room.TimeZone, true)
	if err != nil {
		action.logger.Error("failed to parse time", "error", err)

		msg := "Sorry, I was not able to understand the time."
		go action.storer.SendAndStoreMessage(msg, msg, matrixdb.MessageTypeSetDefaultReminderTimeError, *event)

		return
	}

	loc := time.UTC
	if parsedLoc, err := time.LoadLocation(event.Room.TimeZone); err == nil {
		loc = parsedLoc
	}

	minutesSinceMidnight := uint(timeRemind.In(loc).Hour()*60 + timeRemind.In(loc).Minute())
	event.Channel.DefaultReminderTime = &minutesSinceMidnight

	_, err = action.db.UpdateChannel(event.Channel)
	if err != nil {
		action.logger.Error("failed to update channel", "error", err)

		msg := "Whups, could not save that change. Sorry, try again later."
		go action.storer.SendAndStoreMessage(msg, msg, matrixdb.MessageTypeSetDefaultReminderTimeError, *event)

		return
	}

	msg := fmt.Sprintf("I will use %s as default for your reminders.", format.TimeToHourAndMinute(timeRemind))
	go action.storer.SendAndStoreResponse(msg, matrixdb.MessageTypeSetDefaultReminderTime, *event)
}
