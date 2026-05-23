package daemon_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/require"
)

func TestService_Cleanup(t *testing.T) {
	service, db, outputService := testDaemon(t, false, false, true)

	db.EXPECT().Cleanup(&database.CleanupOpts{
		OlderThan: 365 * 24 * time.Hour,
	}).Return(int64(0), nil)
	outputService.EXPECT().Cleanup().Return(nil)

	go service.Start() //nolint:errcheck

	time.Sleep(time.Millisecond * 5) // give time to execute

	err := service.Stop()
	require.NoError(t, err)
}
