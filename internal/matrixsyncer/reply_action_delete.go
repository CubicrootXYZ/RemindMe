package matrixsyncer

import (
	"fmt"
	"regexp"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

func (s *Syncer) getReplyActionDelete(rtt []database.MessageType) *types.ReplyAction {
	action := &types.ReplyAction{
		Name:         "Delete a reminder",
		Examples:     []string{"delete", "remove", "cancel"},
		Regex:        regexp.MustCompile("(?i)^(delete|remove|cancel)$"),
		ReplyToTypes: rtt,
		Action:       s.replyActionDeleteReminder,
	}
	return action
}

func (s *Syncer) replyActionDeleteReminder(evt *types.MessageEvent, channel *database.Channel, replyMessage *database.Message) error {
	var msg string
	var msgFormatted string
	if replyMessage.ReminderID == nil {
		msg = fmt.Sprintf("Sorry, I could not delete the reminder %d.", replyMessage.ReminderID)
		msgFormatted = msg
		go s.sendAndStoreMessage(asyncmessenger.HTMLMessage(
			msg,
			msgFormatted,
			channel.ChannelIdentifier,
		), channel, database.MessageTypeReminderDeleteFail, 0)

		return errors.ErrIDNotSet
	}
	reminder, err := s.daemon.Database.DeleteReminder(*replyMessage.ReminderID)
	if err != nil {
		msg = fmt.Sprintf("Sorry, I could not delete the reminder %d.", replyMessage.ReminderID)
		msgFormatted = msg
		go s.sendAndStoreMessage(asyncmessenger.HTMLMessage(
			msg,
			msgFormatted,
			channel.ChannelIdentifier,
		), channel, database.MessageTypeReminderDeleteFail, 0)

		return err
	}

	err = s.messenger.DeleteMessageAsync(&asyncmessenger.Delete{
		ExternalIdentifier:        replyMessage.ExternalIdentifier,
		ChannelExternalIdentifier: channel.ChannelIdentifier,
	})
	if err != nil {
		return err
	}

	err = s.messenger.DeleteMessageAsync(&asyncmessenger.Delete{
		ExternalIdentifier:        evt.Event.ID.String(),
		ChannelExternalIdentifier: channel.ChannelIdentifier,
	})
	if err != nil {
		return err
	}

	// Delete all messages regarding this reminder
	messages, err := s.daemon.Database.GetMessagesByReminderID(reminder.ID)
	if err == nil {
		for _, message := range messages {
			err = s.messenger.DeleteMessageAsync(&asyncmessenger.Delete{
				ExternalIdentifier:        message.ExternalIdentifier,
				ChannelExternalIdentifier: channel.ChannelIdentifier,
			})
			if err != nil {
				log.Warn(fmt.Sprintf("Failed to delete message %s with: %s", message.ExternalIdentifier, err.Error()))
			}
		}
	} else {
		log.Warn(fmt.Sprintf("Failed to get messages for reminder %d: %s", reminder.ID, err.Error()))
	}

	_, err = s.daemon.Database.AddMessageFromMatrix(evt.Event.ID.String(), evt.Event.Timestamp, evt.Content, reminder, database.MessageTypeReminderDelete, channel)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to add delete message %s to database: %s", evt.Event.ID.String(), err.Error()))
	}

	msg = fmt.Sprintf("I deleted the reminder \"%s\" (at %s) for you.", reminder.Message, formater.ToLocalTime(reminder.RemindTime, channel.TimeZone))
	msgFormatted = fmt.Sprintf("I <b>deleted</b> the reminder \"%s\" (<i>at %s</i>) for you.", reminder.Message, formater.ToLocalTime(reminder.RemindTime, channel.TimeZone))
	go s.sendAndStoreMessage(asyncmessenger.HTMLMessage(
		msg,
		msgFormatted,
		channel.ChannelIdentifier,
	), channel, database.MessageTypeReminderDeleteSuccess, reminder.ID)

	return nil
}
