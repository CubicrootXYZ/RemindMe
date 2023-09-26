package matrixsyncer

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

func (s *Syncer) getActionList() *types.Action {
	action := &types.Action{
		Name:     "List all reminders",
		Examples: []string{"list", "list reminders", "show", "show reminders", "list my reminders", "reminders"},
		Regex:    regexp.MustCompile("(?i)^((list|show)(| all| the)(| reminders| my reminders)(| please)|^reminders|^reminder)[ ]*$"),
		Action:   s.actionList,
	}
	return action
}

// actionList performs the action "list" that writes all pending reminders to the given channel
func (s *Syncer) actionList(_ *types.MessageEvent, channel *database.Channel) error {
	reminders, err := s.daemon.Database.GetPendingReminders(channel)
	if err != nil {
		return err
	}

	msg := formater.Formater{}

	msg.Title("Open Reminders")
	msg.TextLine("You asked for your open reminders, here they are:")
	msg.NewLine()

	for _, reminder := range reminders {
		msg.BoldLine(reminder.Message)
		msg.ItalicLine("ID " + strconv.FormatUint(uint64(reminder.ID), 10) + " at " + formater.ToLocalTime(reminder.RemindTime, channel.TimeZone) + " " + strings.Join(reminder.GetReminderIcons(), " "))
		msg.NewLine()
	}

	message, messageFormatted := msg.Build()
	go s.sendAndStoreMessage(asyncmessenger.HTMLMessage(
		message,
		messageFormatted,
		channel.ChannelIdentifier,
	), channel, database.MessageTypeReminderList, 0)

	return nil
}
