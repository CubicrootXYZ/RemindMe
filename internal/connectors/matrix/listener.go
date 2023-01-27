package matrix

import (
	"errors"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/encryption"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func (service *service) startListener() error {
	syncer, ok := service.client.Syncer.(*mautrix.DefaultSyncer)
	if !ok {
		return errors.New("syncer of wrong type")
	}

	if service.crypto.enabled {
		olm := encryption.NewOlmMachine(service.client, service.crypto.cryptoStore, service.crypto.stateStore, service.logger.WithField("component", "olm"))

		syncer.OnSync(func(resp *mautrix.RespSync, since string) bool {
			olm.ProcessSyncResponse(resp, since)
			return true
		})
		// TODO syncer.OnEventType(event.EventEncrypted, messageHandler.NewEvent)
		syncer.OnEventType(event.StateEncryption, func(_ mautrix.EventSource, event *event.Event) {
			service.crypto.stateStore.SetEncryptionEvent(event)
		})
	}

	/* TODO syncer.OnEventType(event.EventMessage, messageHandler.NewEvent)
	syncer.OnEventType(event.EventReaction, reactionHandler.NewEvent)
	syncer.OnEventType(event.StateMember, stateMemberHandler.NewEvent)*/

	return service.client.Sync()
}
