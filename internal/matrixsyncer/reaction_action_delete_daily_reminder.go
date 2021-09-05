package matrixsyncer

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getReactionActionDeleteDailyReminder(rat ReactionActionType) *ReactionAction {
	action := &ReactionAction{
		Name:   "Delete the daily message",
		Keys:   []string{"❌"},
		Action: s.reactionActionDeleteDailyReminder,
		Type:   rat,
	}
	return action
}

func (s *Syncer) reactionActionDeleteDailyReminder(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	c, err := s.daemon.Database.UpdateChannel(channel.ID, channel.TimeZone, nil)

	msg := "I will no longer send you a daily message. To reactivate this feature message me with \"set daily reminder at 10:00\"."
	_, err = s.messenger.SendFormattedMessage(msg, msg, c, database.MessageTypeDailyReminderDeleteSuccess, 0)
	return err
}
