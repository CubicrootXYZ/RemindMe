package matrixsyncer

import (
	"fmt"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getReactionActionDelete(rat ReactionActionType) *ReactionAction {
	action := &ReactionAction{
		Name:   "Delete a reminder",
		Keys:   []string{"❌"},
		Action: s.reactionActionDeleteReminder,
		Type:   rat,
	}
	return action
}

func (s *Syncer) reactionActionDeleteReminder(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	var msg string
	var msgFormatted string
	reminder, err := s.daemon.Database.DeleteReminder(message.ReminderID)
	if err != nil {
		msg = fmt.Sprintf("Sorry, I could not delete the reminder %d.", message.ReminderID)
		msgFormatted = msg
		_, err = s.messenger.SendFormattedMessage(msg, msgFormatted, channel.ChannelIdentifier)
		return err
	}

	err = s.messenger.DeleteMessage(message.ExternalIdentifier, channel.ChannelIdentifier)
	if err != nil {
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

	msg = fmt.Sprintf("I deleted the reminder \"%s\" (at %s) for you.", reminder.Message, formater.ToLocalTime(reminder.RemindTime, channel))
	msgFormatted = fmt.Sprintf("I <b>deleted</b> the reminder \"%s\" (<i>at %s</i>) for you.", reminder.Message, formater.ToLocalTime(reminder.RemindTime, channel))
	_, err = s.messenger.SendFormattedMessage(msg, msgFormatted, channel.ChannelIdentifier)
	return err
}
