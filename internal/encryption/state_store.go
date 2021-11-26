package encryption

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"

	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type StateStore struct {
	database types.Database
	config   *configuration.Matrix
}

// NewStateStore returns a new state store
func NewStateStore(database types.Database, config *configuration.Matrix) *StateStore {
	return &StateStore{
		database: database,
		config:   config,
	}
}

// IsEncrypted returns whether a room is encrypted.
func (store *StateStore) IsEncrypted(roomID id.RoomID) bool {
	return store.GetEncryptionEvent(roomID) != nil
}

func (store *StateStore) GetEncryptionEvent(roomID id.RoomID) *event.EncryptionEventContent {
	channels, err := store.database.GetChannelsByChannelIdentifier(roomID.String())
	if err != nil {
		log.Warn(err.Error())
		return nil
	}

	for _, channel := range channels {
		if channel.LastCryptoEvent == "" {
			continue
		}

		encryptionEventJSON := []byte(channel.LastCryptoEvent)

		var encryptionEvent event.EncryptionEventContent
		if err := json.Unmarshal(encryptionEventJSON, &encryptionEvent); err != nil {
			log.Warn(fmt.Sprintf("Failed to unmarshal encryption event JSON: %s. Error: %s", encryptionEventJSON, err))
			return nil
		}
		return &encryptionEvent
	}

	return nil
}

func (store *StateStore) FindSharedRooms(userId id.UserID) []id.RoomID {
	rooms := make([]id.RoomID, 0)
	channels, err := store.database.GetChannelsByUserIdentifier(userId.String())
	if err != nil {
		log.Warn("Could not fetch users rooms: " + err.Error())
		return rooms
	}

	for _, channel := range channels {
		if channel.LastCryptoEvent != "" {
			rooms = append(rooms, id.RoomID(channel.ChannelIdentifier))
		}
	}

	return rooms
}

func (store *StateStore) SetMembership(event *event.Event) {
	// Do not do anything, this is already handled elsewhere
}

func (store *StateStore) SetEncryptionEvent(event *event.Event) {
	if event == nil {
		return
	}

	channels, err := store.database.GetChannelsByChannelIdentifier(event.RoomID.String())
	if err != nil {
		log.Warn("Failed setting encryption event: " + err.Error())
		return
	}

	var encryptionEventJson []byte
	encryptionEventJson, err = json.Marshal(event)
	if err != nil {
		encryptionEventJson = nil
	}

	for _, channel := range channels {
		channel.LastCryptoEvent = string(encryptionEventJson)
		if err := store.database.ChannelSaveChanges(&channel); err != nil {
			log.Warn("Failed saving encryption event: " + err.Error())
		}
	}
}

func (store *StateStore) GetUserIDs(roomID string) []id.UserID {
	userIDs := make([]id.UserID, 0)
	userIDs = append(userIDs, id.UserID(fmt.Sprintf("@%s:%s", store.config.Username, strings.ReplaceAll(strings.ReplaceAll(store.config.Homeserver, "https://", ""), "http://", ""))))

	channels, err := store.database.GetChannelsByChannelIdentifier(roomID)
	if err != nil {
		log.Warn("Failed getting rooms: " + err.Error())
	}

	for _, channel := range channels {
		userIDs = append(userIDs, id.UserID(channel.UserIdentifier))
	}

	return userIDs
}
