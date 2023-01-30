package database_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testMessage() *database.MatrixMessage {
	user, err := service.NewUser(testUser())
	if err != nil {
		panic(err)
	}

	messageID := "abcde"
	for err == nil {
		messageID = fmt.Sprintf("abcde%d", rand.Int()) //nolint:gosec
		_, err = service.GetMessageByID(messageID)
	}

	return &database.MatrixMessage{
		ID:            messageID,
		User:          *user,
		Room:          user.Rooms[0],
		Body:          "TEST",
		BodyFormatted: "<b>test</b>",
		SendAt:        time.Now(),
		Type:          "testtype",
		Incoming:      true,
	}
}

func TestService_GetMessageByID(t *testing.T) {
	messageBefore, err := service.NewMessage(testMessage())
	require.NoError(t, err)

	messageAfter, err := service.GetMessageByID(messageBefore.ID)
	require.NoError(t, err)

	assertMessagesEqual(t, messageBefore, messageAfter)
}

func TestService_GetLastMessage(t *testing.T) {
	_, err := service.NewMessage(testMessage())
	require.NoError(t, err)
	time.Sleep(time.Millisecond * 50)

	_, err = service.NewMessage(testMessage())
	require.NoError(t, err)
	time.Sleep(time.Second)

	messageBefore, err := service.NewMessage(testMessage())
	require.NoError(t, err)

	messageAfter, err := service.GetLastMessage()
	require.NoError(t, err)

	assertMessagesEqual(t, messageBefore, messageAfter)
}

func TestService_DeleteAllMessagesFromRoom(t *testing.T) {
	message, err := service.NewMessage(testMessage())
	require.NoError(t, err)

	err = service.DeleteAllMessagesFromRoom(message.RoomID)
	require.NoError(t, err)

	_, err = service.GetMessageByID(message.ID)
	assert.ErrorIs(t, err, database.ErrNotFound)
}

func assertMessagesEqual(t *testing.T, a *database.MatrixMessage, b *database.MatrixMessage) {
	t.Helper()

	// We do not load them, so ignore them
	a.User.Rooms = nil
	b.User.Rooms = nil
	a.Room.Users = nil
	b.Room.Users = nil

	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.Body, b.Body)
	assert.Equal(t, a.BodyFormatted, b.BodyFormatted)
	assert.Equal(t, a.SendAt.UTC().Format(time.RFC3339), b.SendAt.UTC().Format(time.RFC3339))
	assert.Equal(t, a.Type, b.Type)
	assert.Equal(t, a.Incoming, b.Incoming)
	assertUsersEqual(t, &a.User, &b.User)
	assertRoomsEqual(t, &a.Room, &b.Room)
}
