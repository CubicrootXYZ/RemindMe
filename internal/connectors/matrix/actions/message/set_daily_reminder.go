package message

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

var setDailyReminderRegex = regexp.MustCompile("(?i)^(set|update|change|)[ ]*(the|a|my|)[ ]*(daily reminder|daily info|daily message).*")

// SetDailyReminderAction sets the time for the daily reminder.
type SetDailyReminderAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
	storer    *msghelper.Storer
}

// Configure is called on startup and sets all dependencies.
func (action *SetDailyReminderAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
	action.storer = msghelper.NewStorer(matrixDB, messenger, logger)
}

// Name of the action
func (action *SetDailyReminderAction) Name() string {
	return "Set Daily Reminder"
}

// GetDocu returns the documentation for the action.
func (action *SetDailyReminderAction) GetDocu() (title, explaination string, examples []string) {
	return "Set Daily Reminder",
		"Set the time to send a reminder of todays events",
		[]string{"daily reminder at 9am", "daily reminder at 13:00", "set the daily info at 4:00"}
}

// Selector defines a regex on what messages the action should be used.
func (action *SetDailyReminderAction) Selector() *regexp.Regexp {
	return setDailyReminderRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *SetDailyReminderAction) HandleEvent(event *matrix.MessageEvent) {
	dbMsg := mapping.MessageFromEvent(event)
	dbMsg.Type = matrixdb.MessageTypeSetDailyReminder
	_, err := action.matrixDB.NewMessage(dbMsg)
	if err != nil {
		action.logger.Errorf("Could not save message: %s", err.Error())
	}

	timeRemind, err := format.ParseTime(event.Content.Body, event.Room.TimeZone, true)
	if err != nil {
		action.logger.Err(err)
		msg := "Sorry, I was not able to understand the time."
		go action.storer.SendAndStoreMessage(msg, msg, matrixdb.MessageTypeSetDailyReminderError, *event)
		return
	}

	minutesSinceMidnight := uint(timeRemind.Hour()*60 + timeRemind.Minute())

	event.Channel.DailyReminder = &minutesSinceMidnight
	_, err = action.db.UpdateChannel(event.Channel)
	if err != nil {
		action.logger.Err(err)
		msg := "Whups, could not save that change. Sorry, try again later."
		go action.storer.SendAndStoreMessage(msg, msg, matrixdb.MessageTypeSetDailyReminderError, *event)
		return
	}

	msg := fmt.Sprintf("I will send you a daily overview at %s. To disable the reminder message me with \"delete daily reminder\".", format.TimeToHourAndMinute(timeRemind))
	go action.storer.SendAndStoreResponse(msg, matrixdb.MessageTypeSetDailyReminder, *event)
}
