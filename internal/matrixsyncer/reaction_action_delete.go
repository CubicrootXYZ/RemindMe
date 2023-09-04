package matrixsyncer

import (
	"fmt"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getReactionActionDelete(rat types.ReactionActionType) *types.ReactionAction {
	action := &types.ReactionAction{
		Name:   "Delete a reminder",
		Keys:   []string{"❌"},
		Action: s.reactionActionDeleteReminder,
		Type:   rat,
	}
	return action
}

func (s *Syncer) reactionActionDeleteReminder(message *database.Message, _ *event.ReactionEventContent, _ *event.Event, channel *database.Channel) error {
	if channel == nil {
		return errors.ErrEmptyChannel
	}
	var msg string
	var msgFormatted string
	if message.ReminderID == nil {
		msg = fmt.Sprintf("Sorry, I could not delete the reminder %d.", message.ReminderID)
		msgFormatted = msg
		go s.sendAndStoreMessage(asyncmessenger.HTMLMessage(
			msg,
			msgFormatted,
			channel.ChannelIdentifier,
		), channel, database.MessageTypeReminderDeleteFail, 0)
		return errors.ErrIDNotSet
	}
	reminder, err := s.daemon.Database.DeleteReminder(*message.ReminderID)
	if err != nil || reminder == nil {
		msg = fmt.Sprintf("Sorry, I could not delete the reminder %d.", message.ReminderID)
		go s.sendAndStoreMessage(asyncmessenger.PlainTextMessage(
			msg,
			channel.ChannelIdentifier,
		), channel, database.MessageTypeReminderDeleteFail, *message.ReminderID)
		return err
	}

	err = s.messenger.DeleteMessageAsync(&asyncmessenger.Delete{
		ExternalIdentifier:        message.ExternalIdentifier,
		ChannelExternalIdentifier: channel.ChannelIdentifier,
	})
	if err != nil {
		log.Info("Could not delete message, are you sure the bot has the permission to do so? " + err.Error())
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

	msg = fmt.Sprintf("I deleted the reminder \"%s\" (at %s) for you.", reminder.Message, formater.ToLocalTime(reminder.RemindTime, channel.TimeZone))
	msgFormatted = fmt.Sprintf("I <b>deleted</b> the reminder \"%s\" (<i>at %s</i>) for you.", reminder.Message, formater.ToLocalTime(reminder.RemindTime, channel.TimeZone))
	go s.sendAndStoreMessage(asyncmessenger.HTMLMessage(
		msg,
		msgFormatted,
		channel.ChannelIdentifier,
	), channel, database.MessageTypeReminderDeleteSuccess, reminder.ID)

	return nil
}
