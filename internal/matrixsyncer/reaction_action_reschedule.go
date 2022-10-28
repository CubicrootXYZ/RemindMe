package matrixsyncer

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
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

	respondToEvent := &types.MessageEvent{
		Event: &event.Event{
			Sender: id.UserID(message.Channel.UserIdentifier),
			ID:     id.EventID(message.ExternalIdentifier),
		},
		Content: &event.MessageEventContent{
			FormattedBody: message.BodyHTML,
			Body:          message.Body,
		},
	}

	requestMessage, err := s.daemon.Database.GetLastMessageByTypeForReminder(database.MessageTypeReminderRequest, *message.ReminderID)
	if err == nil {
		respondToEvent = &types.MessageEvent{
			Event: &event.Event{
				Sender: id.UserID(requestMessage.Channel.UserIdentifier),
				ID:     id.EventID(requestMessage.ExternalIdentifier),
			},
			Content: &event.MessageEventContent{
				FormattedBody: requestMessage.BodyHTML,
				Body:          requestMessage.Body,
			},
		}
	}

	reminder, err := s.daemon.Database.GetReminderForChannelIDByID(channel.ChannelIdentifier, int(*message.ReminderID))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	newRemindTime := reminder.RemindTime.Add(24 * time.Hour)
	for time.Until(newRemindTime) < 1*time.Hour {
		newRemindTime = newRemindTime.Add(24 * time.Hour)
	}

	_, err = s.daemon.Database.UpdateReminder(reminder.ID, newRemindTime, reminder.RepeatInterval, reminder.RepeatMax)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = s.messenger.DeleteMessageAsync(&asyncmessenger.Delete{
		ExternalIdentifier:        message.ExternalIdentifier,
		ChannelExternalIdentifier: channel.ChannelIdentifier,
	})
	if err != nil {
		log.Info("Could not delete message, are you sure the bot has the permission to do so? " + err.Error())
	}

	msg := "Rescheduled that reminder to tomorrow."
	err = s.messenger.SendResponseAsync(asyncmessenger.PlainTextResponse(
		msg,
		respondToEvent.Event.ID.String(),
		respondToEvent.Content.Body,
		respondToEvent.Event.Sender.String(),
		channel.ChannelIdentifier,
	))

	return nil
}
