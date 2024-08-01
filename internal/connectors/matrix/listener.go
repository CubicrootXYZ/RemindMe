package matrix

import (
	"errors"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func (service *service) startListener() error {
	syncer, ok := service.client.Syncer.(*mautrix.DefaultSyncer)
	if !ok {
		return errors.New("syncer of wrong type")
	}

	syncer.OnEventType(event.EventMessage, service.MessageEventHandler)
	syncer.OnEventType(event.EventReaction, service.ReactionEventHandler)
	syncer.OnEventType(event.StateMember, service.EventStateHandler)

	return service.client.Sync()
}
