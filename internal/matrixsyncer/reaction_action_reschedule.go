package matrixsyncer

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func (s *Syncer) getReactionActionReschedule(rat types.ReactionActionType) *types.ReactionAction {
	action := &types.ReactionAction{
		Name:   "Reschedule reminder to tomorrow",
		Keys:   []string{"🔄"},
		Action: s.reactionActionRescheduleReminder,
		Type:   rat,
	}
	return action
}

func (s *Syncer) reactionActionRescheduleReminder(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	if channel == nil {
		return errors.ErrEmptyChannel
	}

	if message.ReminderID == nil {
		return ErrMessageHasNoReminder
	}

	reminder, err := s.daemon.Database.GetReminderForChannelIDByID(channel.ChannelIdentifier, int(*message.ReminderID))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	newRemindTime := reminder.RemindTime
	for newRemindTime.Sub(reminder.RemindTime) < 24*time.Hour {
		newRemindTime = newRemindTime.Add(24 * time.Hour)
	}

	_, err = s.daemon.Database.UpdateReminder(reminder.ID, newRemindTime, reminder.RepeatInterval, reminder.RepeatMax)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = s.messenger.DeleteMessage(message.ExternalIdentifier, channel.ChannelIdentifier)
	if err != nil {
		log.Info("Could not delete message, are you sure the bot has the permission to do so? " + err.Error())
	}

	msg := "Rescheduled that reminder to tomorrow."
	_, err = s.messenger.SendReplyToEvent(msg, &types.MessageEvent{
		Event: &event.Event{
			Sender: id.UserID(message.Channel.UserIdentifier),
			ID:     id.EventID(message.ExternalIdentifier),
		},
		Content: &event.MessageEventContent{
			FormattedBody: message.BodyHTML,
			Body:          message.Body,
		},
	}, channel, database.MessageTypeReminderDeleteSuccess)
	return err
}
