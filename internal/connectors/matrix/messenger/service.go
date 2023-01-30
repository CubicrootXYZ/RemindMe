package asyncmessenger

import (
	"sync"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/encryption"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type service struct {
	roomUserCache roomCache
	config        *Config
	client        MatrixClient
	db            database.Service
	logger        gologger.Logger
	state         *state
}

type Config struct {
	crypto *cryptoTools
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

func NewMessenger(config *Config, db database.Service, matrixClient MatrixClient, logger gologger.Logger) (Messenger, error) {
	return &service{
		roomUserCache: make(roomCache),
		config:        config,
		client:        matrixClient,
		db:            db,
		logger:        logger,
		state: &state{
			rateLimitedUntilMutex: sync.Mutex{},
		},
	}, nil
}

// sendMessageEvent sends a message event to matrix, will take care of encryption if available
func (messenger *service) sendMessageEvent(messageEvent *messageEvent, roomID string, eventType event.Type) (*mautrix.RespSendEvent, error) {
	if messenger.config.crypto.stateStore != nil && eventType != event.EventReaction {
		if messenger.config.crypto.stateStore.IsEncrypted(id.RoomID(roomID)) && messenger.config.crypto.olm != nil {
			resp, err := messenger.sendMessageEventEncrypted(messageEvent, roomID, eventType)
			if err == nil {
				return resp, nil
			}
		}
	}

	messenger.logger.Infof("Sending message to room %s", roomID)
	return messenger.client.SendMessageEvent(id.RoomID(roomID), eventType, &messageEvent)
}

func (messenger *service) sendMessageEventEncrypted(messageEvent *messageEvent, roomID string, eventType event.Type) (*mautrix.RespSendEvent, error) {
	messenger.config.crypto.cryptoMutex.Lock()

	encrypted, err := messenger.config.crypto.olm.EncryptMegolmEvent(id.RoomID(roomID), eventType, messageEvent)

	if err == crypto.SessionExpired || err == crypto.SessionNotShared || err == crypto.NoGroupSession {
		err = messenger.config.crypto.olm.ShareGroupSession(id.RoomID(roomID), messenger.getUserIDsInRoom(id.RoomID(roomID)))
		if err != nil {
			messenger.config.crypto.cryptoMutex.Unlock()
			return nil, err
		}

		encrypted, err = messenger.config.crypto.olm.EncryptMegolmEvent(id.RoomID(roomID), eventType, messageEvent)
	}
	messenger.config.crypto.cryptoMutex.Unlock()
	if err != nil {
		return nil, err
	}

	enrichCleartext(encrypted, messageEvent)

	messenger.logger.Infof("Sending encrypted message to room %s", roomID)
	return messenger.client.SendMessageEvent(id.RoomID(roomID), event.EventEncrypted, encrypted)
}

// enrichCleartext adds parts of the encrypted event back into the cleartext event as specified by the matrix spec
func enrichCleartext(encryptedEvent *event.EncryptedEventContent, evt *messageEvent) {
	if evt.RelatesTo.EventID == "" && evt.RelatesTo.InReplyTo == nil {
		return
	}

	encryptedEvent.RelatesTo = &event.RelatesTo{}
	encryptedEvent.RelatesTo.EventID = id.EventID(evt.RelatesTo.EventID)
	encryptedEvent.RelatesTo.Key = evt.RelatesTo.Key
	encryptedEvent.RelatesTo.Type = event.RelationType(evt.RelatesTo.RelType)

	if evt.RelatesTo.InReplyTo != nil {
		encryptedEvent.RelatesTo.InReplyTo = &event.InReplyTo{
			EventID: id.EventID(evt.RelatesTo.InReplyTo.EventID),
		}
	}
}

func (messenger *service) getUserIDsInRoom(roomID id.RoomID) []id.UserID {
	// Check cache first
	if users := messenger.roomUserCache.GetUsers(roomID); users != nil {
		return users
	}

	userIDs := make([]id.UserID, 0)
	members, err := messenger.client.JoinedMembers(roomID)
	if err != nil {
		messenger.logger.Err(err)
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

func (messenger *service) encounteredRateLimit() {
	messenger.state.rateLimitedUntilMutex.Lock()
	messenger.state.rateLimitedUntil = time.Now().Add(time.Minute)
	messenger.state.rateLimitedUntilMutex.Unlock()
}
