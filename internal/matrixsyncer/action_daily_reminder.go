package matrixsyncer

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getActionSetDailyReminder() *Action {
	action := &Action{
		Name:     "Set daily reminder time",
		Examples: []string{"daily reminder at 9am", "daily reminder at 13:00", "set the daily info at 4:00"},
		Regex:    "(?i)^(set|update|change|)[ ]*(the|a|my|)[ ]*(daily reminder|daily info|daily message).*",
		Action:   s.actionSetDailyReminder,
	}
	return action
}

// actionSetDailyReminder sets the daily reminder time
func (s *Syncer) actionSetDailyReminder(evt *event.Event, channel *database.Channel) error {
	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		log.Warn("Event is not a message event. Can not handle it")
		return errors.ErrMatrixEventWrongType
	}

	_, err := s.daemon.Database.AddMessageFromMatrix(evt.ID.String(), time.Now().Unix(), content, nil, database.MessageTypeDailyReminderUpdate, channel)
	if err != nil {
		log.Error("Could not save message: " + err.Error())
	}

	timeRemind, err := formater.ParseTime(content.Body, channel, true)

	minutesSinceMidnight := uint(timeRemind.Hour()*60 + timeRemind.Minute())

	c, err := s.daemon.Database.UpdateChannel(channel.ID, channel.TimeZone, &minutesSinceMidnight)
	if err != nil {
		_, err = s.messenger.SendReplyToEvent("Sorry, I was not able to save that.", evt, channel, database.MessageTypeDailyReminderUpdateFail)
		return err
	}

	_, err = s.messenger.SendReplyToEvent(fmt.Sprintf("I will send you a daily overview at %s", formater.TimeToHourAndMinute(timeRemind)), evt, c, database.MessageTypeDailyReminderUpdateSuccess)
	return err
}

func (s *Syncer) getActionDeleteDailyReminder() *Action {
	action := &Action{
		Name:     "Delete daily reminder time",
		Examples: []string{"remove daily reminder", "delete daily message"},
		Regex:    "(?i)^(remove|delete|cancel)[ ]*(the|a|my|)[ ]*(daily reminder|daily info|daily message).*",
		Action:   s.actionDeleteDailyReminder,
	}
	return action
}

// actionDeleteDailyReminder deletes the daily reminder
func (s *Syncer) actionDeleteDailyReminder(evt *event.Event, channel *database.Channel) error {
	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		log.Warn("Event is not a message event. Can not handle it")
		return errors.ErrMatrixEventWrongType
	}

	_, err := s.daemon.Database.AddMessageFromMatrix(evt.ID.String(), time.Now().Unix(), content, nil, database.MessageTypeDailyReminderDelete, channel)
	if err != nil {
		log.Error("Could not save message: " + err.Error())
	}

	c, err := s.daemon.Database.UpdateChannel(channel.ID, channel.TimeZone, nil)
	if err != nil {
		_, err = s.messenger.SendReplyToEvent("Sorry, I was not able to save that.", evt, channel, database.MessageTypeDailyReminderDeleteFail)
		return err
	}

	_, err = s.messenger.SendReplyToEvent("I will no longer send you a daily message.", evt, c, database.MessageTypeDailyReminderDeleteSuccess)
	return err
}
