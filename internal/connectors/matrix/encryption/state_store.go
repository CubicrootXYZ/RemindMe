package encryption

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// StateStore holds encryption states.
type StateStore struct {
	database database.Service
	config   *StateStoreConfig
	logger   gologger.Logger
}

type StateStoreConfig struct {
	Username   string
	Homeserver string
}

// NewStateStore returns a new state store
func NewStateStore(database database.Service, config *StateStoreConfig, logger gologger.Logger) *StateStore {
	return &StateStore{
		database: database,
		config:   config,
		logger:   logger,
	}
}

// IsEncrypted returns whether a room is encrypted.
func (store *StateStore) IsEncrypted(roomID id.RoomID) bool {
	return store.GetEncryptionEvent(roomID) != nil
}

func (store *StateStore) GetEncryptionEvent(roomID id.RoomID) *event.EncryptionEventContent {
	room, err := store.database.GetRoomByRoomID(roomID.String())
	if err != nil {
		store.logger.Err(err)
		return nil
	}

	if room.LastCryptoEvent == "" {
		return nil
	}

	encryptionEventJSON := []byte(room.LastCryptoEvent)

	var encryptionEvent event.EncryptionEventContent
	if err := json.Unmarshal(encryptionEventJSON, &encryptionEvent); err != nil {
		store.logger.Errorf("Failed to unmarshal encryption event JSON: %s. Error: %s", encryptionEventJSON, err)
		return nil
	}
	return &encryptionEvent
}

func (store *StateStore) FindSharedRooms(userID id.UserID) []id.RoomID {
	rooms := make([]id.RoomID, 0)
	user, err := store.database.GetUserByID(userID.String())
	if err != nil {
		store.logger.Errorf("Could not fetch users rooms: " + err.Error())
		return rooms
	}

	for _, room := range user.Rooms {
		if room.LastCryptoEvent != "" {
			rooms = append(rooms, id.RoomID(room.RoomID))
		}
	}

	return rooms
}

func (store *StateStore) SetMembership(_ *event.Event) {
	// Do not do anything, this is already handled elsewhere
}

func (store *StateStore) SetEncryptionEvent(event *event.Event) {
	if event == nil {
		return
	}

	room, err := store.database.GetRoomByRoomID(event.RoomID.String())
	if err != nil {
		store.logger.Errorf("Failed setting encryption event: " + err.Error())
		return
	}

	var encryptionEventJSON []byte
	encryptionEventJSON, err = json.Marshal(event)
	if err != nil {
		encryptionEventJSON = nil
	}
	room.LastCryptoEvent = string(encryptionEventJSON)
	if _, err := store.database.UpdateRoom(room); err != nil {
		store.logger.Errorf("Failed saving encryption event: " + err.Error())
	}
}

func (store *StateStore) GetUserIDs(roomID string) []id.UserID {
	userIDs := make([]id.UserID, 0)
	userIDs = append(userIDs, id.UserID(fmt.Sprintf("@%s:%s", store.config.Username, strings.ReplaceAll(strings.ReplaceAll(store.config.Homeserver, "https://", ""), "http://", ""))))

	room, err := store.database.GetRoomByRoomID(roomID)
	if err != nil {
		store.logger.Errorf("Failed getting rooms: " + err.Error())
		return userIDs
	}

	for _, user := range room.Users {
		userIDs = append(userIDs, id.UserID(user.ID))
	}

	return userIDs
}
