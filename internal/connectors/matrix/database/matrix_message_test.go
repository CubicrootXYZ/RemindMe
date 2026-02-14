package database_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testMessage() *matrixdb.MatrixMessage {
	user, err := service.NewUser(testUser())
	if err != nil {
		panic(err)
	}

	messageID := "abcde"
	for err == nil {
		messageID = fmt.Sprintf("abcde%d", rand.Int()) //nolint:gosec
		_, err = service.GetMessageByID(messageID)
	}

	return &matrixdb.MatrixMessage{
		ID:            messageID,
		User:          user,
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

func TestService_GetMessageByIDWithNotFoundError(t *testing.T) {
	_, err := service.GetMessageByID("dslkfajgö9w40tßegpjdfapg")
	require.ErrorIs(t, err, matrixdb.ErrNotFound)
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
	assert.ErrorIs(t, err, matrixdb.ErrNotFound)
}

func TestService_GetEventMessageByOutputAndEvent(t *testing.T) {
	c := database.Channel{}
	err := gormDB.Save(&c).Error
	require.NoError(t, err)

	evt := &database.Event{
		Channel: c,
		Time:    time.Now(),
	}
	err = gormDB.Save(evt).Error
	require.NoError(t, err)

	messageBefore := testMessage()
	messageBefore.EventID = &evt.ID
	messageBefore.Type = matrixdb.MessageTypeNewEvent
	messageBefore.Incoming = true
	messageBefore, err = service.NewMessage(messageBefore)
	require.NoError(t, err)

	messageAfter, err := service.GetEventMessageByOutputAndEvent(evt.ID, messageBefore.RoomID, matrix.OutputType)
	require.NoError(t, err)
	assertMessagesEqual(t, messageBefore, messageAfter)
}

func TestService_GetEventMessageByOutputAndEventWithRoomNotFound(t *testing.T) {
	c := database.Channel{}
	err := gormDB.Save(&c).Error
	require.NoError(t, err)

	evt := &database.Event{
		Channel: c,
		Time:    time.Now(),
	}
	err = gormDB.Save(evt).Error
	require.NoError(t, err)

	messageBefore := testMessage()
	messageBefore.EventID = &evt.ID
	messageBefore.Type = matrixdb.MessageTypeNewEvent
	messageBefore.Incoming = true
	messageBefore, err = service.NewMessage(messageBefore)
	require.NoError(t, err)

	_, err = service.GetEventMessageByOutputAndEvent(evt.ID, messageBefore.RoomID+59384, matrix.OutputType)
	assert.ErrorIs(t, err, matrixdb.ErrNotFound)
}

func TestService_GetEventMessageByOutputAndEventWithMessageNotFound(t *testing.T) {
	c := database.Channel{}
	err := gormDB.Save(&c).Error
	require.NoError(t, err)

	evt := &database.Event{
		Channel: c,
		Time:    time.Now(),
	}
	err = gormDB.Save(evt).Error
	require.NoError(t, err)

	messageBefore := testMessage()
	messageBefore.EventID = &evt.ID
	messageBefore.Type = matrixdb.MessageTypeNewEvent
	messageBefore.Incoming = false
	messageBefore, err = service.NewMessage(messageBefore)
	require.NoError(t, err)

	_, err = service.GetEventMessageByOutputAndEvent(evt.ID, messageBefore.RoomID, matrix.OutputType)
	assert.ErrorIs(t, err, matrixdb.ErrNotFound)
}

func TestService_ListMessages(t *testing.T) {
	msg1, err := service.NewMessage(testMessage())
	require.NoError(t, err)
	time.Sleep(time.Millisecond * 50)

	msg2 := testMessage()
	msg2.Type = matrixdb.MessageTypeChangeEvent
	msg2, err = service.NewMessage(msg2)
	require.NoError(t, err)
	time.Sleep(time.Second)

	// List first message.
	msgs, err := service.ListMessages(matrixdb.ListMessageOpts{
		RoomID:   &msg1.RoomID,
		Incoming: new(true),
	})
	require.NoError(t, err)
	require.Len(t, msgs, 1)
	assertMessagesEqual(t, msg1, &msgs[0])

	// List second message.
	msgs, err = service.ListMessages(matrixdb.ListMessageOpts{
		RoomID: &msg2.RoomID,
		Type:   &matrixdb.MessageTypeChangeEvent,
	})
	require.NoError(t, err)
	require.Len(t, msgs, 1)
	assertMessagesEqual(t, msg2, &msgs[0])
}

func assertMessagesEqual(t *testing.T, a *matrixdb.MatrixMessage, b *matrixdb.MatrixMessage) {
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
	assertUsersEqual(t, a.User, b.User)
	assertRoomsEqual(t, &a.Room, &b.Room)
}
