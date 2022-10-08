package asyncmessenger

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

	return messageEvent
}

func (messenger *messenger) SendReactionAsync(reaction *Reaction) error {
	go messenger.sendMessage(reaction.toEvent(), reaction.ChannelExternalIdentifier, 10, time.Second*15)

	return nil
}