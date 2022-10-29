package matrixsyncer

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getReactionActionDeleteDailyReminder(rat types.ReactionActionType) *types.ReactionAction {
	action := &types.ReactionAction{
		Name:   "Delete the daily message",
		Keys:   []string{"❌"},
		Action: s.reactionActionDeleteDailyReminder,
		Type:   rat,
	}
	return action
}

func (s *Syncer) reactionActionDeleteDailyReminder(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	c, err := s.daemon.Database.UpdateChannel(channel.ID, channel.TimeZone, nil, channel.Role)
	if err != nil {
		log.Error(err.Error())
		go s.sendAndStoreMessage(asyncmessenger.PlainTextMessage(
			"Sorry I was not able to delete the daily reminder.",
			c.ChannelIdentifier,
		), c, database.MessageTypeDailyReminderDeleteFail, 0)
		return err
	}

	go s.sendAndStoreMessage(asyncmessenger.PlainTextMessage(
		"I will no longer send you a daily message. To reactivate this feature message me with \"set daily reminder at 10:00\".",
		c.ChannelIdentifier,
	), c, database.MessageTypeDailyReminderDeleteSuccess, 0)

	return nil
}
