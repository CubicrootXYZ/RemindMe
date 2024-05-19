package database_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testEvent() *matrixdb.MatrixEvent {
	user, err := service.NewUser(testUser())
	if err != nil {
		panic(err)
	}

	eventID := "abcde"
	for err == nil {
		eventID = fmt.Sprintf("abcde%d", rand.Int()) //nolint:gosec
		_, err = service.GetEventByID(eventID)
	}

	return &matrixdb.MatrixEvent{
		ID:     eventID,
		User:   *user,
		Room:   user.Rooms[0],
		SendAt: time.Now(),
		Type:   "testtype",
	}
}

func TestService_GetEventByID(t *testing.T) {
	eventBefore, err := service.NewEvent(testEvent())
	require.NoError(t, err)

	eventAfter, err := service.GetEventByID(eventBefore.ID)
	require.NoError(t, err)

	assertEventsEqual(t, eventBefore, eventAfter)
}

func TestService_GetLastEvent(t *testing.T) {
	_, err := service.NewEvent(testEvent())
	require.NoError(t, err)
	time.Sleep(time.Millisecond * 50)

	_, err = service.NewEvent(testEvent())
	require.NoError(t, err)
	time.Sleep(time.Second)

	eventBefore, err := service.NewEvent(testEvent())
	require.NoError(t, err)

	eventAfter, err := service.GetLastEvent()
	require.NoError(t, err)

	assertEventsEqual(t, eventBefore, eventAfter)
}

func TestService_DeleteAllEventsFromRoom(t *testing.T) {
	event, err := service.NewEvent(testEvent())
	require.NoError(t, err)

	err = service.DeleteAllEventsFromRoom(event.Room.ID)
	require.NoError(t, err)

	_, err = service.GetEventByID(event.Type)
	assert.ErrorIs(t, err, matrixdb.ErrNotFound)
}

func assertEventsEqual(t *testing.T, a *matrixdb.MatrixEvent, b *matrixdb.MatrixEvent) {
	t.Helper()

	// We do not load them, so ignore them
	a.User.Rooms = nil
	b.User.Rooms = nil
	a.Room.Users = nil
	b.Room.Users = nil

	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.SendAt.UTC().Format(time.RFC3339), b.SendAt.UTC().Format(time.RFC3339))
	assert.Equal(t, a.Type, b.Type)
	assertUsersEqual(t, &a.User, &b.User)
	assertRoomsEqual(t, &a.Room, &b.Room)
}
