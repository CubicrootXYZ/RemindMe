package messenger

import "time"

type Reaction struct {
	Reaction                  string
	MessageExternalIdentifier string
	ChannelExternalIdentifier string
}

func (reaction *Reaction) toEvent() *messageEvent {
	messageEvent := &messageEvent{
		Type: messageTypeReaction,
	}
	messageEvent.RelatesTo.EventID = reaction.MessageExternalIdentifier
	messageEvent.RelatesTo.Key = reaction.Reaction
	messageEvent.RelatesTo.RelType = relationAnnotiation
	messageEvent.RelatesTo.InReplyTo = nil

	return messageEvent
}

func (messenger *service) SendReactionAsync(reaction *Reaction) error {
	if messenger.config.DisableReactions {
		return nil
	}

	go func() {
		_, _ = messenger.sendMessage(reaction.toEvent(), reaction.ChannelExternalIdentifier, 10, time.Second*15)
	}()

	return nil
}
