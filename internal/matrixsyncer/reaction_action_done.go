package matrixsyncer

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
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

	err = s.messenger.DeleteMessage(message.ExternalIdentifier, channel.ChannelIdentifier)
	if err != nil {
		log.Info("Could not delete message, are you sure the bot has the permission to do so? " + err.Error())
	}

	msg := "Great work! I marked that reminder as done."
	_, err = s.messenger.SendReplyToEvent(msg, respondToEvent, channel, database.MessageTypeReminderDeleteSuccess)
	return err
}
