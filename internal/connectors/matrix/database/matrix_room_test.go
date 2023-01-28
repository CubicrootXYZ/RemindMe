package database_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRoom() *database.MatrixRoom {
	roomID := "!12345678:example.org"
	var err error
	for err == nil {
		roomID = fmt.Sprintf("!%d:example.org", rand.Int()) //nolint:gosec
		_, err = service.GetRoomByID(roomID)
	}

	return &database.MatrixRoom{
		RoomID:          roomID,
		LastCryptoEvent: "{}",
	}
}

func TestGetRoomByID(t *testing.T) {
	roomBefore, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	roomAfter, err := service.GetRoomByID(roomBefore.RoomID)
	require.NoError(t, err)

	assertRoomsEqual(t, roomBefore, roomAfter)
}

func TestGetRoomByIDWithRoomNotFound(t *testing.T) {
	_, err := service.GetRoomByID("abc")
	assert.ErrorIs(t, err, database.ErrNotFound)
}

func TestUpdateRoom(t *testing.T) {
	roomBefore, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	roomBefore.LastCryptoEvent = "{\"a\": \"b\"}"
	_, err = service.UpdateRoom(roomBefore)
	require.NoError(t, err)

	roomAfter, err := service.GetRoomByID(roomBefore.RoomID)
	require.NoError(t, err)
	assertRoomsEqual(t, roomBefore, roomAfter)
}

func assertRoomsEqual(t *testing.T, a *database.MatrixRoom, b *database.MatrixRoom) {
	t.Helper()

	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.CreatedAt.UTC(), b.CreatedAt.UTC())
	assert.Equal(t, a.UpdatedAt.UTC(), b.UpdatedAt.UTC())
	if !a.DeletedAt.Valid {
		assert.False(t, b.DeletedAt.Valid)
	} else {
		assert.Equal(t, a.DeletedAt.Time.UTC(), b.DeletedAt.Time.UTC())
	}
	assert.Equal(t, a.RoomID, b.RoomID)
	if len(a.Users) == 0 {
		assert.Equal(t, 0, len(b.Users))
	} else {
		assert.Equal(t, a.Users, b.Users)
	}

	assert.Equal(t, a.LastCryptoEvent, b.LastCryptoEvent)
}
