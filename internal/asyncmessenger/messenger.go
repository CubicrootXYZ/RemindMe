package asyncmessenger

import (
	"fmt"
	"sync"
	"time"

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
	roomUserCache roomCache
	config        *configuration.Config
	client        MatrixClient
	db            types.Database
	debug         bool
	cryptoTools   *cryptoTools
	state         *state
}

type cryptoTools struct {
	olm         *crypto.OlmMachine
	stateStore  *encryption.StateStore
	cryptoMutex sync.Mutex // Since the crypto foo relies on a single sqlite, only one process at a time is allowed to access it
}

type state struct {
	rateLimitedUntil      time.Time // If we run into a rate limit this will tell us to stop operation
	rateLimitedUntilMutex sync.Mutex
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
		roomUserCache: make(roomCache),
		config:        config,
		client:        matrixClient,
		db:            db,
		debug:         debug,
		state: &state{
			rateLimitedUntilMutex: sync.Mutex{},
		},
		cryptoTools: &cryptoTools{
			olm:         olm,
			stateStore:  stateStore,
			cryptoMutex: sync.Mutex{},
		},
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
		err = messenger.cryptoTools.olm.ShareGroupSession(id.RoomID(roomID), messenger.getUserIDsInRoom(id.RoomID(roomID)))
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

func (messenger *messenger) getUserIDsInRoom(roomID id.RoomID) []id.UserID {
	// Check cache first
	if users := messenger.roomUserCache.GetUsers(roomID); users != nil {
		return users
	}

	userIDs := make([]id.UserID, 0)
	members, err := messenger.client.JoinedMembers(roomID)
	if err != nil {
		log.Warn(err.Error())
		return userIDs
	}

	i := 0
	for userID := range members.Joined {
		userIDs = append(userIDs, userID)
		i++
	}

	messenger.roomUserCache.AddUsers(roomID, userIDs)
	return userIDs
}

func (messenger *messenger) encounteredRateLimit() {
	messenger.state.rateLimitedUntilMutex.Lock()
	messenger.state.rateLimitedUntil = time.Now().Add(time.Minute)
	messenger.state.rateLimitedUntilMutex.Unlock()
}
