package matrixsyncer

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getReactionsAddTime(rat types.ReactionActionType) []*types.ReactionAction {
	actions := make([]*types.ReactionAction, 0)
	actions = append(actions, &types.ReactionAction{
		Name:   "Add 1 hour",
		Keys:   []string{"1️⃣"},
		Action: s.reactionActionAdd1Hour,
		Type:   rat,
	})
	actions = append(actions, &types.ReactionAction{
		Name:   "Add 2 hours",
		Keys:   []string{"2️⃣"},
		Action: s.reactionActionAdd2Hours,
		Type:   rat,
	})
	actions = append(actions, &types.ReactionAction{
		Name:   "Add 3 hours",
		Keys:   []string{"3️⃣"},
		Action: s.reactionActionAdd3Hours,
		Type:   rat,
	})
	actions = append(actions, &types.ReactionAction{
		Name:   "Add 4 hours",
		Keys:   []string{"4️⃣"},
		Action: s.reactionActionAdd4Hours,
		Type:   rat,
	})
	actions = append(actions, &types.ReactionAction{
		Name:   "Add 5 hours",
		Keys:   []string{"5️⃣"},
		Action: s.reactionActionAdd5Hours,
		Type:   rat,
	})
	actions = append(actions, &types.ReactionAction{
		Name:   "Add 6 hours",
		Keys:   []string{"6️⃣"},
		Action: s.reactionActionAdd6Hours,
		Type:   rat,
	})
	actions = append(actions, &types.ReactionAction{
		Name:   "Add 7 hours",
		Keys:   []string{"7️⃣"},
		Action: s.reactionActionAdd7Hours,
		Type:   rat,
	})
	actions = append(actions, &types.ReactionAction{
		Name:   "Add 8 hours",
		Keys:   []string{"8️⃣"},
		Action: s.reactionActionAdd8Hours,
		Type:   rat,
	})
	actions = append(actions, &types.ReactionAction{
		Name:   "Add 9 hours",
		Keys:   []string{"9️⃣"},
		Action: s.reactionActionAdd9Hours,
		Type:   rat,
	})
	actions = append(actions, &types.ReactionAction{
		Name:   "Add 10 hours",
		Keys:   []string{"🔟"},
		Action: s.reactionActionAdd10Hours,
		Type:   rat,
	})
	actions = append(actions, &types.ReactionAction{
		Name:   "Add 1 day",
		Keys:   []string{"➕"},
		Action: s.reactionActionAdd1Day,
		Type:   rat,
	})

	return actions
}

func (s *Syncer) reactionActionAdd1Hour(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	return s.reactionActionAddXHours(message, content, evt, channel, 1*time.Hour)
}

func (s *Syncer) reactionActionAdd2Hours(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	return s.reactionActionAddXHours(message, content, evt, channel, 2*time.Hour)
}

func (s *Syncer) reactionActionAdd3Hours(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	return s.reactionActionAddXHours(message, content, evt, channel, 3*time.Hour)
}

func (s *Syncer) reactionActionAdd4Hours(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	return s.reactionActionAddXHours(message, content, evt, channel, 4*time.Hour)
}

func (s *Syncer) reactionActionAdd5Hours(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	return s.reactionActionAddXHours(message, content, evt, channel, 5*time.Hour)
}

func (s *Syncer) reactionActionAdd6Hours(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	return s.reactionActionAddXHours(message, content, evt, channel, 6*time.Hour)
}

func (s *Syncer) reactionActionAdd7Hours(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	return s.reactionActionAddXHours(message, content, evt, channel, 7*time.Hour)
}

func (s *Syncer) reactionActionAdd8Hours(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	return s.reactionActionAddXHours(message, content, evt, channel, 8*time.Hour)
}

func (s *Syncer) reactionActionAdd9Hours(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	return s.reactionActionAddXHours(message, content, evt, channel, 9*time.Hour)
}

func (s *Syncer) reactionActionAdd10Hours(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	return s.reactionActionAddXHours(message, content, evt, channel, 10*time.Hour)
}

func (s *Syncer) reactionActionAdd1Day(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	return s.reactionActionAddXHours(message, content, evt, channel, 24*time.Hour)
}

// reactionAddXHours to referentiate all other actions here
func (s *Syncer) reactionActionAddXHours(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel, duration time.Duration) error {
	log.Debug(fmt.Sprintf("Adding %d minutes, with reaction %s (event %s)", duration/time.Minute, content.RelatesTo.Key, evt.ID))

	if message.ReminderID == nil {
		msg := fmt.Sprintf("Sorry, I could not delete the reminder %d.", message.ReminderID)
		msgFormatted := msg
		go s.sendAndStoreMessage(asyncmessenger.HTMLMessage(
			msg,
			msgFormatted,
			channel.ChannelIdentifier,
		), channel, database.MessageTypeReminderRecurringFail, 0)

		return errors.ErrIDNotSet
	}

	reminder, err := s.daemon.Database.UpdateReminder(*message.ReminderID, addTimeOrFromNow(message.Reminder.RemindTime, duration), 0, 0)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Reminder \"%s\" rescheduled to %s", reminder.Message, formater.ToLocalTime(reminder.RemindTime, channel.TimeZone))

	go s.sendAndStoreMessage(asyncmessenger.PlainTextMessage(
		msg,
		channel.ChannelIdentifier,
	), channel, database.MessageTypeReminderUpdateSuccess, reminder.ID)

	return nil
}

// addTimeOrFromNow If baseTime is in future add the duration to the basetime otherwise add it to the current time
func addTimeOrFromNow(baseTime time.Time, duration time.Duration) time.Time {
	if now := time.Now(); baseTime.Sub(now) < 0 {
		return now.Add(duration)
	}

	return baseTime.Add(duration)
}
