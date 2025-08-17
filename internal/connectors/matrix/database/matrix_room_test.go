package database_test

import (
	"fmt"
	"math/rand"
	"testing"

	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRoom() *matrixdb.MatrixRoom {
	roomID := "!12345678:example.org"

	var err error
	for err == nil {
		roomID = fmt.Sprintf("!%d:example.org", rand.Int()) //nolint:gosec
		_, err = service.GetRoomByRoomID(roomID)
	}

	return &matrixdb.MatrixRoom{
		RoomID: roomID,
	}
}

func TestService_ListInputRoomsByChannel(t *testing.T) {
	roomBefore, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	channel := &database.Channel{}
	require.NoError(t, gormDB.Save(channel).Error)
	input := &database.Input{
		InputType: "matrix",
		InputID:   roomBefore.ID,
		ChannelID: channel.ID,
	}
	require.NoError(t, gormDB.Save(input).Error)

	rooms, err := service.ListInputRoomsByChannel(channel.ID)
	require.NoError(t, err)

	foundRoom := false

	for _, r := range rooms {
		if r.ID == roomBefore.ID {
			foundRoom = true
			break
		}
	}

	assert.True(t, foundRoom, "room is not in response")
}

func TestService_ListInputRoomsByChannelWithEmpty(t *testing.T) {
	roomBefore, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	channel := &database.Channel{}
	require.NoError(t, gormDB.Save(channel).Error)
	output := &database.Output{
		OutputType: "matrix",
		OutputID:   roomBefore.ID,
		ChannelID:  channel.ID,
	}
	require.NoError(t, gormDB.Save(output).Error)

	rooms, err := service.ListInputRoomsByChannel(channel.ID)
	require.NoError(t, err)

	for _, r := range rooms {
		if r.ID == roomBefore.ID {
			assert.Fail(t, "room should not be in response")
		}
	}
}

func TestService_ListOutputRoomsByChannel(t *testing.T) {
	roomBefore, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	channel := &database.Channel{}
	require.NoError(t, gormDB.Save(channel).Error)
	input := &database.Output{
		OutputType: "matrix",
		OutputID:   roomBefore.ID,
		ChannelID:  channel.ID,
	}
	require.NoError(t, gormDB.Save(input).Error)

	rooms, err := service.ListOutputRoomsByChannel(channel.ID)
	require.NoError(t, err)

	foundRoom := false

	for _, r := range rooms {
		if r.ID == roomBefore.ID {
			foundRoom = true
			break
		}
	}

	assert.True(t, foundRoom, "room is not in response")
}

func TestService_ListOutputRoomsByChannelWithEmpty(t *testing.T) {
	roomBefore, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	channel := &database.Channel{}
	require.NoError(t, gormDB.Save(channel).Error)
	output := &database.Input{
		InputType: "matrix",
		InputID:   roomBefore.ID,
		ChannelID: channel.ID,
	}
	require.NoError(t, gormDB.Save(output).Error)

	rooms, err := service.ListOutputRoomsByChannel(channel.ID)
	require.NoError(t, err)

	for _, r := range rooms {
		if r.ID == roomBefore.ID {
			assert.Fail(t, "room should not be in response")
		}
	}
}

func TestService_GetRoomByID(t *testing.T) {
	roomBefore, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	roomAfter, err := service.GetRoomByID(roomBefore.ID)
	require.NoError(t, err)

	assertRoomsEqual(t, roomBefore, roomAfter)
}

func TestGetRoomByIDWithRoomNotFound(t *testing.T) {
	_, err := service.GetRoomByID(9999)
	assert.ErrorIs(t, err, matrixdb.ErrNotFound)
}

func TestService_GetRoomByRoomID(t *testing.T) {
	roomBefore, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	roomAfter, err := service.GetRoomByRoomID(roomBefore.RoomID)
	require.NoError(t, err)

	assertRoomsEqual(t, roomBefore, roomAfter)
}

func TestGetRoomByRoomIDWithRoomNotFound(t *testing.T) {
	_, err := service.GetRoomByRoomID("abc")
	assert.ErrorIs(t, err, matrixdb.ErrNotFound)
}

func TestService_UpdateRoom(t *testing.T) {
	roomBefore, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	_, err = service.UpdateRoom(roomBefore)
	require.NoError(t, err)

	roomAfter, err := service.GetRoomByRoomID(roomBefore.RoomID)
	require.NoError(t, err)
	assertRoomsEqual(t, roomBefore, roomAfter)
}

func TestService_DeleteRoom(t *testing.T) {
	user, room := createRoomWithUser(t)

	roomAfter, err := service.GetRoomByID(room.ID)
	require.NoError(t, err)
	require.Len(t, roomAfter.Users, 1)

	err = service.DeleteRoom(room.ID)
	require.NoError(t, err)

	_, err = service.GetRoomByRoomID(room.RoomID)
	require.ErrorIs(t, err, matrixdb.ErrNotFound)

	userAfter, err := service.GetUserByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, userAfter.ID)
}

func TestService_GetRoomCount(t *testing.T) {
	_, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	cnt, err := service.GetRoomCount()
	require.NoError(t, err)

	assert.LessOrEqual(t, int64(1), cnt)
}

func TestService_AddUserToRoom(t *testing.T) {
	user := testUser()
	roomBefore, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	roomBefore, err = service.AddUserToRoom(user.ID, roomBefore)
	require.NoError(t, err)

	roomAfter, err := service.GetRoomByID(roomBefore.ID)
	require.NoError(t, err)

	assertRoomsEqual(t, roomBefore, roomAfter)

	userFound := false

	for _, u := range roomAfter.Users {
		if u.ID == user.ID {
			userFound = true
			break
		}
	}

	assert.Truef(t, userFound, "user '%s' was not in room altough added", user.ID)
}

func TestService_AddUserToRoomWithUserAlreadyExists(t *testing.T) {
	user, err := service.NewUser(testUser())
	require.NoError(t, err)

	roomBefore, err := service.NewRoom(testRoom())
	require.NoError(t, err)

	roomBefore, err = service.AddUserToRoom(user.ID, roomBefore)
	require.NoError(t, err)

	roomBefore.Users[0].Rooms = nil // GetRoomByID does not set rooms

	roomAfter, err := service.GetRoomByID(roomBefore.ID)
	require.NoError(t, err)

	assertRoomsEqual(t, roomBefore, roomAfter)

	userFound := false

	for _, u := range roomAfter.Users {
		if u.ID == user.ID {
			userFound = true
			break
		}
	}

	assert.Truef(t, userFound, "user '%s' was not in room altough added", user.ID)
}

func assertRoomsEqual(t *testing.T, a *matrixdb.MatrixRoom, b *matrixdb.MatrixRoom) {
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
		assert.Empty(t, b.Users)
	} else {
		assert.Equal(t, a.Users, b.Users)
	}
}
