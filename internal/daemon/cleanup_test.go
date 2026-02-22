package daemon_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestService_Cleanup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, db, _ := testDaemon(ctrl, false, false, true)

	db.EXPECT().Cleanup(&database.CleanupOpts{
		OlderThan: 365 * 24 * time.Hour,
	}).Return(int64(0), nil).MinTimes(1)

	go service.Start() //nolint:errcheck

	time.Sleep(time.Millisecond * 5) // give time to execute

	err := service.Stop()
	require.NoError(t, err)
}
