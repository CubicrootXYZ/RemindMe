package matrixsyncer

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"gorm.io/gorm"
)

func (s *Syncer) getActionDeleteReminder() *types.Action {
	action := &types.Action{
		Name:     "Delete a reminder by ID",
		Examples: []string{"delete reminder 1", "remove 68"},
		Regex:    regexp.MustCompile("(?i)(^(delete|remove)[ ]*(reminder|)[ ]+[0-9]+)$"),
		Action:   s.actionDeleteReminder,
	}
	return action
}

func (s *Syncer) actionDeleteReminder(evt *types.MessageEvent, channel *database.Channel) error {
	reminderID, err := formater.GetSuffixInt(evt.Content.Body)
	if err != nil {
		log.Error(err.Error())
		msg := "Whupsy, I expected a number in that message but could not find it."
		_, err = s.messenger.SendReplyToEvent(msg, evt, channel, database.MessageTypeDoNotSave)
		return err
	}

	reminder, err := s.daemon.Database.GetReminderForChannelIDByID(channel.ChannelIdentifier, reminderID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := "Sorry, I do not know this reminder."
		_, err = s.messenger.SendReplyToEvent(msg, evt, channel, database.MessageTypeDoNotSave)
		return err
	} else if err != nil {
		log.Error(err.Error())
		msg := "Oh no, this did not work, sorry."
		_, err = s.messenger.SendReplyToEvent(msg, evt, channel, database.MessageTypeDoNotSave)
		return err
	}

	_, err = s.daemon.Database.DeleteReminder(reminder.ID)
	if err != nil {
		log.Error(err.Error())
		msg := "Sorry, this did not work."
		_, err = s.messenger.SendReplyToEvent(msg, evt, channel, database.MessageTypeDoNotSave)
		return err
	}

	// Delete all messages regarding this reminder
	messages, err := s.daemon.Database.GetMessagesByReminderID(reminder.ID)
	if err == nil {
		for _, message := range messages {
			err = s.messenger.DeleteMessage(message.ExternalIdentifier, channel.ChannelIdentifier)
			if err != nil {
				log.Warn(fmt.Sprintf("Failed to delete message %s with: %s", message.ExternalIdentifier, err.Error()))
			}
		}
	} else {
		log.Warn(fmt.Sprintf("Failed to get messages for reminder %d: %s", reminder.ID, err.Error()))
	}

	msgFormater := formater.Formater{}
	msgFormater.TextLine("Deleted the reminder: ")
	msgFormater.QuoteLine(reminder.Message)
	msg, formattedMsg := msgFormater.Build()

	_, err = s.messenger.SendFormattedMessage(msg, formattedMsg, channel, database.MessageTypeReminderDelete, reminder.ID)
	return err
}
