package asyncmessenger

import (
	"fmt"
	"sync"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/encryption"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type messenger struct {
	config      *configuration.Config
	client      *mautrix.Client
	db          types.Database
	debug       bool
	cryptoTools *cryptoTools
}

type cryptoTools struct {
	olm         *crypto.OlmMachine
	stateStore  *encryption.StateStore
	cryptoMutex sync.Mutex // Since the crypto foo relies on a single sqlite, only one process at a time is allowed to access it
}

func NewMessenger(debug bool, config *configuration.Config, db types.Database, cryptoStore crypto.Store, stateStore *encryption.StateStore, matrixClient *mautrix.Client) (Messenger, error) {
	var olm *crypto.OlmMachine
	if config.MatrixBotAccount.E2EE {
		olm = encryption.GetOlmMachine(debug, matrixClient, cryptoStore, db, stateStore)
		olm.AllowUnverifiedDevices = true
		olm.ShareKeysToUnverifiedDevices = true
		err := olm.Load()
		if err != nil {
			return nil, err
		}
	}

	return &messenger{
		config:      config,
		client:      matrixClient,
		db:          db,
		olm:         olm,
		stateStore:  stateStore,
		debug:       debug,
		cryptoMutex: sync.Mutex{},
	}, nil
}

// sendMessageEvent sends a message event to matrix, will take care of encryption if available
func (messenger *messenger) sendMessageEvent(messageEvent *messageEvent, roomID string, eventType event.Type) (*mautrix.RespSendEvent, error) {
	if messenger.cryptoTools.stateStore != nil && eventType != event.EventReaction {
		if messenger.cryptoTools.stateStore.IsEncrypted(id.RoomID(roomID)) && messenger.cryptoTools.olm != nil {
			resp, err := messenger.sendMessageEventEncrypted(messageEvent, roomID, eventType)
			if err == nil {
				return resp, nil
			}
		}
	}

	log.Info(fmt.Sprintf("Sending message to room %s", roomID))
	return messenger.client.SendMessageEvent(id.RoomID(roomID), eventType, &messageEvent)
}

func (messenger *messenger) sendMessageEventEncrypted(messageEvent *messageEvent, roomID string, eventType event.Type) (*mautrix.RespSendEvent, error) {
	messenger.cryptoTools.cryptoMutex.Lock()

	encrypted, err := messenger.cryptoTools.olm.EncryptMegolmEvent(id.RoomID(roomID), eventType, messageEvent)

	if err == crypto.SessionExpired || err == crypto.SessionNotShared || err == crypto.NoGroupSession {
		err = messenger.cryptoTools.olm.ShareGroupSession(id.RoomID(roomID), messenger.getUserIDs(id.RoomID(roomID)))
		if err != nil {
			messenger.cryptoTools.cryptoMutex.Unlock()
			return nil, err
		}

		encrypted, err = messenger.cryptoTools.olm.EncryptMegolmEvent(id.RoomID(roomID), eventType, messageEvent)
	}
	messenger.cryptoTools.cryptoMutex.Unlock()
	if err != nil {
		return nil, err
	}

	log.Info(fmt.Sprintf("Sending encrypted message to room %s", roomID))
	return messenger.client.SendMessageEvent(id.RoomID(roomID), event.EventEncrypted, encrypted)
}
