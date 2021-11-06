package encryption

import (
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type stateStore struct {
}

// NewStateStore returns a new state store
func NewStateStore() crypto.StateStore {
	return &stateStore{}
}

// TODO get from database?
func (store *stateStore) IsEncrypted(roomID id.RoomID) bool {
	return true
}

func (store *stateStore) GetEncryptionEvent(roomID id.RoomID) *event.EncryptionEventContent {
	return &event.EncryptionEventContent{}
}

func (store *stateStore) FindSharedRooms(userID id.UserID) []id.RoomID {
	return []id.RoomID{}
}
