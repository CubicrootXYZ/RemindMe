package database_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Cleanup(t *testing.T) {
	channel, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	oldEvent, err := service.NewEvent(&database.Event{
		Time:      time.Now().Add(-48 * time.Hour).UTC(),
		Duration:  time.Hour,
		Message:   "cleanup-old",
		Active:    true,
		ChannelID: channel.ID,
	})
	require.NoError(t, err)

	newEvent, err := service.NewEvent(&database.Event{
		Time:      time.Now().UTC(),
		Duration:  time.Hour,
		Message:   "cleanup-new",
		Active:    true,
		ChannelID: channel.ID,
	})
	require.NoError(t, err)

	deleted, err := service.Cleanup(&database.CleanupOpts{
		OlderThan: 24 * time.Hour,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, deleted, int64(1))

	var events []database.Event

	err = gormDB.Where("id in(?)", []uint{oldEvent.ID, newEvent.ID}).Find(&events).Error
	require.NoError(t, err)

	require.Len(t, events, 1)
	assert.Equal(t, newEvent.Message, events[0].Message)
}
