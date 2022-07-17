package matrixsyncer

import (
	"fmt"
	"regexp"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

func (s *Syncer) getActionSetDailyReminder() *types.Action {
	action := &types.Action{
		Name:     "Set daily reminder time",
		Examples: []string{"daily reminder at 9am", "daily reminder at 13:00", "set the daily info at 4:00"},
		Regex:    regexp.MustCompile("(?i)^(set|update|change|)[ ]*(the|a|my|)[ ]*(daily reminder|daily info|daily message).*"),
		Action:   s.actionSetDailyReminder,
	}
	return action
}

// actionSetDailyReminder sets the daily reminder time
func (s *Syncer) actionSetDailyReminder(evt *types.MessageEvent, channel *database.Channel) error {
	_, err := s.daemon.Database.AddMessageFromMatrix(evt.Event.ID.String(), time.Now().Unix(), evt.Content, nil, database.MessageTypeDailyReminderUpdate, channel)
	if err != nil {
		log.Error("Could not save message: " + err.Error())
	}

	timeRemind, err := formater.ParseTime(evt.Content.Body, channel, true)
	if err != nil {
		_, err = s.messenger.SendReplyToEvent("Sorry, I was not able to understand the time.", evt, channel, database.MessageTypeDailyReminderUpdateFail)
		return err
	}

	minutesSinceMidnight := uint(timeRemind.Hour()*60 + timeRemind.Minute())

	c, err := s.daemon.Database.UpdateChannel(channel.ID, channel.TimeZone, &minutesSinceMidnight, channel.Role)
	if err != nil {
		_, err = s.messenger.SendReplyToEvent("Sorry, I was not able to save that.", evt, channel, database.MessageTypeDailyReminderUpdateFail)
		return err
	}

	_, err = s.messenger.SendReplyToEvent(fmt.Sprintf("I will send you a daily overview at %s. To disable this message me with \"delete daily reminder\".", formater.TimeToHourAndMinute(timeRemind)), evt, c, database.MessageTypeDailyReminderUpdateSuccess)
	return err
}

func (s *Syncer) getActionDeleteDailyReminder() *types.Action {
	action := &types.Action{
		Name:     "Delete daily reminder time",
		Examples: []string{"remove daily reminder", "delete daily message"},
		Regex:    regexp.MustCompile("(?i)^(remove|delete|cancel)[ ]*(the|a|my|)[ ]*(daily reminder|daily info|daily message).*"),
		Action:   s.actionDeleteDailyReminder,
	}
	return action
}

// actionDeleteDailyReminder deletes the daily reminder
func (s *Syncer) actionDeleteDailyReminder(evt *types.MessageEvent, channel *database.Channel) error {
	_, err := s.daemon.Database.AddMessageFromMatrix(evt.Event.ID.String(), time.Now().Unix(), evt.Content, nil, database.MessageTypeDailyReminderDelete, channel)
	if err != nil {
		log.Error("Could not save message: " + err.Error())
	}

	c, err := s.daemon.Database.UpdateChannel(channel.ID, channel.TimeZone, nil, channel.Role)
	if err != nil {
		_, err = s.messenger.SendReplyToEvent("Sorry, I was not able to save that.", evt, channel, database.MessageTypeDailyReminderDeleteFail)
		return err
	}

	_, err = s.messenger.SendReplyToEvent("I will no longer send you a daily message. To reactivate this feature message me with \"set daily reminder at 10:00\".", evt, c, database.MessageTypeDailyReminderDeleteSuccess)
	return err
}
