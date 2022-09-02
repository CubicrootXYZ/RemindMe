package matrixsyncer

import (
	"regexp"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

var regexChangeReminder = regexp.MustCompile("((change|update|set)[ ]+(reminder|reminder id|)[ ]*[0-9]+)")

func (s *Syncer) getActionChangeReminder() *types.Action {
	action := &types.Action{
		Name:     "Change a reminder by ID",
		Examples: []string{"change reminder 1 to tomorrow", "update 68 to Saturday 4 pm"},
		Regex:    regexp.MustCompile("(?i)(^(change|update|set)[ ]+(reminder|reminder id|)[ ]*[0-9]+)"),
		Action:   s.actionChangeReminder,
	}
	return action
}

func (s *Syncer) actionChangeReminder(evt *types.MessageEvent, channel *database.Channel) error {
	match := regexChangeReminder.Find([]byte(evt.Content.Body))
	if match == nil {
		msg := "Ups, seems like there is a reminder ID missing in your message."
		_, err := s.messenger.SendReplyToEvent(msg, evt, channel, database.MessageTypeDoNotSave)
		return err
	}

	reminderID, err := formater.GetSuffixInt(string(match))
	if err != nil {
		msg := "Ups, seems like there is a reminder ID missing in your message."
		_, err := s.messenger.SendReplyToEvent(msg, evt, channel, database.MessageTypeDoNotSave)
		return err
	}

	newTime, err := formater.ParseTime(evt.Content.Body, channel, false)
	if err != nil {
		msg := "Ehm, sorry to say that, but I was not able to understand the time to schedule the reminder to."
		_, err := s.messenger.SendReplyToEvent(msg, evt, channel, database.MessageTypeDoNotSave)
		return err
	}

	reminder, err := s.daemon.Database.GetReminderForChannelIDByID(channel.ChannelIdentifier, reminderID)
	if err != nil {
		msg := "This reminder is not in my database."
		_, err := s.messenger.SendReplyToEvent(msg, evt, channel, database.MessageTypeDoNotSave)
		return err
	}

	reminder.RemindTime = newTime
	_, err = s.daemon.Database.UpdateReminder(reminder.ID, newTime, reminder.RepeatInterval, reminder.RepeatMax)
	if err != nil {
		log.Error(err.Error())
		msg := "Whups, this did not work, sorry."
		_, err = s.messenger.SendReplyToEvent(msg, evt, channel, database.MessageTypeDoNotSave)
		return err
	}

	msgFormater := formater.Formater{}
	msgFormater.TextLine("I rescheduled your reminder")
	msgFormater.QuoteLine(reminder.Message)
	msgFormater.Text("to ")
	msgFormater.Text(formater.ToLocalTime(newTime, channel.TimeZone))

	msg, formattedMsg := msgFormater.Build()
	_, err = s.messenger.SendFormattedMessage(msg, formattedMsg, channel, database.MessageTypeReminderUpdate, reminder.ID)

	return err
}
