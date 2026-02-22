package database_test

import (
	"testing"
	"time"

	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/require"
)

func TestService_Cleanup(t *testing.T) {
	// Event for the message with event_id set (must be kept even when old)
	var c database.Channel

	err := gormDB.Save(&c).Error
	require.NoError(t, err)

	evt := &database.Event{ChannelID: c.ID, Time: time.Now()}
	err = gormDB.Save(evt).Error
	require.NoError(t, err)

	// Old message without event_id -> should be deleted
	oldMsg := testMessage()
	oldMsg.SendAt = time.Now().Add(-200 * 24 * time.Hour)
	oldMsg, err = service.NewMessage(oldMsg)
	require.NoError(t, err)

	// New message -> should be kept
	newMsg := testMessage()
	newMsg.SendAt = time.Now()
	newMsg, err = service.NewMessage(newMsg)
	require.NoError(t, err)

	// Old message with event_id set -> should be kept
	oldWithEvent := testMessage()
	oldWithEvent.SendAt = time.Now().Add(-200 * 24 * time.Hour)
	oldWithEvent.EventID = &evt.ID
	oldWithEvent, err = service.NewMessage(oldWithEvent)
	require.NoError(t, err)

	err = service.Cleanup()
	require.NoError(t, err)

	_, err = service.GetMessageByID(oldMsg.ID)
	require.ErrorIs(t, err, matrixdb.ErrNotFound)

	_, err = service.GetMessageByID(newMsg.ID)
	require.NoError(t, err)

	_, err = service.GetMessageByID(oldWithEvent.ID)
	require.NoError(t, err)
}
