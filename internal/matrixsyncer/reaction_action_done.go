package matrixsyncer

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/random"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getReactionActionDone(rat types.ReactionActionType) *types.ReactionAction {
	action := &types.ReactionAction{
		Name:   "Mark reminder as done",
		Keys:   []string{"✅"},
		Action: s.reactionActionDoneReminder,
		Type:   rat,
	}
	return action
}

func (s *Syncer) reactionActionDoneReminder(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	if channel == nil {
		return errors.ErrEmptyChannel
	}

	if message.ReminderID == nil {
		return ErrMessageHasNoReminder
	}

	requestMessage, err := s.daemon.Database.GetLastMessageByTypeForReminder(database.MessageTypeReminderRequest, *message.ReminderID)
	if err == nil {
		err = s.messenger.SendResponseAsync(asyncmessenger.PlainTextResponse(
			"I marked that reminder as done. "+random.MotivationalSentence(),
			requestMessage.ExternalIdentifier,
			requestMessage.Body,
			requestMessage.Channel.UserIdentifier,
			channel.ChannelIdentifier,
		))
		if err != nil {
			log.Error("Failed sending response: " + err.Error())
		}
	}

	err = s.messenger.DeleteMessageAsync(&asyncmessenger.Delete{
		ExternalIdentifier:        message.ExternalIdentifier,
		ChannelExternalIdentifier: channel.ChannelIdentifier,
	})
	if err != nil {
		log.Info("Could not delete message, are you sure the bot has the permission to do so? " + err.Error())
	}

	return nil
}
