package encryption_test

import (
	"errors"
	"testing"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database/mocks"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/encryption"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"maunium.net/go/mautrix/id"
)

func stateStore(ctrl *gomock.Controller) (*encryption.StateStore, *mocks.MockService) {
	db := mocks.NewMockService(ctrl)

	return encryption.NewStateStore(db, &encryption.StateStoreConfig{
		Username:   "bot",
		Homeserver: "example.com",
	}, gologger.New(gologger.LogLevelDebug, 0)), db
}

func TestStore_IsEncrypted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store, db := stateStore(ctrl)

	db.EXPECT().GetRoomByID("abcd").Return(
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

	db.EXPECT().GetRoomByID("abcd").Return(
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

	db.EXPECT().GetRoomByID("abcd").Return(
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

	db.EXPECT().GetRoomByID("abcd").Return(
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

	db.EXPECT().GetRoomByID("abcd").Return(
		&database.MatrixRoom{
			LastCryptoEvent: `{"}`,
		},
		nil,
	)

	content := store.GetEncryptionEvent(id.RoomID("abcd"))
	assert.Nil(t, content)
}
