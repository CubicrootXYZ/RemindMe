package encryption_test

import (
	"errors"
	"testing"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/encryption"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func stateStore(ctrl *gomock.Controller) (*encryption.StateStore, *matrixdb.MockService) {
	db := matrixdb.NewMockService(ctrl)

	return encryption.NewStateStore(db, &encryption.StateStoreConfig{
		Username:   "bot",
		Homeserver: "example.com",
	}, gologger.New(gologger.LogLevelDebug, 0)), db
}

func TestStore_IsEncrypted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetRoomByRoomID("abcd").Return(
		&database.MatrixRoom{
			LastCryptoEvent: `{"algorithm":"my algo", "rotation_period_ms": 12}`,
		},
		nil,
	)

	assert.True(t, store.IsEncrypted("abcd"))
}

func TestStore_IsEncryptedWithIsNot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetRoomByRoomID("abcd").Return(
		&database.MatrixRoom{
			LastCryptoEvent: ``,
		},
		nil,
	)

	assert.False(t, store.IsEncrypted("abcd"))
}

func TestStore_GetEncryptionEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetRoomByRoomID("abcd").Return(
		&database.MatrixRoom{
			LastCryptoEvent: `{"algorithm":"my algo", "rotation_period_ms": 12}`,
		},
		nil,
	)

	content := store.GetEncryptionEvent(id.RoomID("abcd"))
	require.NotNil(t, content)

	assert.Equal(t, "my algo", string(content.Algorithm))
	assert.Equal(t, int64(12), content.RotationPeriodMillis)
}

func TestStore_GetEncryptionEventWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetRoomByRoomID("abcd").Return(
		nil,
		errors.New("test"),
	)

	content := store.GetEncryptionEvent(id.RoomID("abcd"))
	assert.Nil(t, content)
}

func TestStore_GetEncryptionEventWithInvalidEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetRoomByRoomID("abcd").Return(
		&database.MatrixRoom{
			LastCryptoEvent: `{"}`,
		},
		nil,
	)

	content := store.GetEncryptionEvent(id.RoomID("abcd"))
	assert.Nil(t, content)
}

func TestStore_FindSharedRooms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetUserByID("abcd").Return(&database.MatrixUser{
		Rooms: []database.MatrixRoom{
			{
				RoomID:          "123",
				LastCryptoEvent: "{}",
			},
			{
				RoomID:          "456",
				LastCryptoEvent: "{}",
			},
			{
				RoomID: "789",
			},
		},
	}, nil)

	roomIDs := store.FindSharedRooms("abcd")
	require.Equal(t, 2, len(roomIDs), "expected 2 room IDs")
	assert.Equal(t, "123", roomIDs[0].String())
	assert.Equal(t, "456", roomIDs[1].String())
}

func TestStore_FindSharedRoomsWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetUserByID("abcd").Return(nil, errors.New("test"))

	roomIDs := store.FindSharedRooms("abcd")
	require.Equal(t, 0, len(roomIDs), "expected 0 room IDs")
}

func TestStore_SetEncryptionEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetRoomByRoomID("abcd").Return(&database.MatrixRoom{
		RoomID: "abcd",
	}, nil)
	db.EXPECT().UpdateRoom(&database.MatrixRoom{
		RoomID:          "abcd",
		LastCryptoEvent: `{"type":"","room_id":"abcd","content":{}}`,
	}).Return(nil, nil)

	event := &event.Event{
		RoomID: id.RoomID("abcd"),
	}
	store.SetEncryptionEvent(event)
}

func TestStore_SetEncryptionEventWithGetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetRoomByRoomID("abcd").Return(&database.MatrixRoom{
		RoomID: "abcd",
	}, errors.New("test"))

	event := &event.Event{
		RoomID: id.RoomID("abcd"),
	}
	store.SetEncryptionEvent(event)
}

func TestStore_SetEncryptionEventWithUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetRoomByRoomID("abcd").Return(&database.MatrixRoom{
		RoomID: "abcd",
	}, nil)
	db.EXPECT().UpdateRoom(&database.MatrixRoom{
		RoomID:          "abcd",
		LastCryptoEvent: `{"type":"","room_id":"abcd","content":{}}`,
	}).Return(nil, errors.New("test"))

	event := &event.Event{
		RoomID: id.RoomID("abcd"),
	}
	store.SetEncryptionEvent(event)
}

func TestStore_GetUserIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetRoomByRoomID("abcd").Return(&database.MatrixRoom{
		Users: []database.MatrixUser{
			{
				ID: "123",
			},
			{
				ID: "456",
			},
		},
	}, nil)

	userIDs := store.GetUserIDs("abcd")

	require.Equal(t, 3, len(userIDs), "expected 3 user IDs")
	assert.Equal(t, "@bot:example.com", userIDs[0].String())
	assert.Equal(t, "123", userIDs[1].String())
	assert.Equal(t, "456", userIDs[2].String())
}

func TestStore_GetUserIDsWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetRoomByRoomID("abcd").Return(nil, errors.New("test"))

	userIDs := store.GetUserIDs("abcd")

	require.Equal(t, 1, len(userIDs), "expected 1 user IDs")
	assert.Equal(t, "@bot:example.com", userIDs[0].String())
}
