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

	if service.crypto.enabled {
		syncer.OnSync(func(resp *mautrix.RespSync, since string) bool {
			service.crypto.olm.ProcessSyncResponse(resp, since) // TODO this is panicing
			return true
		})
		syncer.OnEventType(event.EventEncrypted, service.MessageEventHandler)
		syncer.OnEventType(event.StateEncryption, func(_ mautrix.EventSource, event *event.Event) {
			service.crypto.stateStore.SetEncryptionEvent(event)
		})
	}

	syncer.OnEventType(event.EventMessage, service.MessageEventHandler)
	syncer.OnEventType(event.EventReaction, service.ReactionEventHandler)
	syncer.OnEventType(event.StateMember, service.EventStateHandler)

	return service.client.Sync()
}
