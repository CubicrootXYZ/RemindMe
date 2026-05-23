package daemon_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func testDatabaseChannel() database.Channel {
	dailyReminder := uint(0)

	return database.Channel{
		DailyReminder: &dailyReminder,
		Outputs:       []database.Output{*testDatabaseOutput()},
	}
}

func TestService_SendOutDailyReminders(t *testing.T) {
	service, db, outputService := testDaemon(t, false, true, false)

	channel := testDatabaseChannel()
	output := testOutput()

	db.EXPECT().GetChannels().Return([]database.Channel{channel}, nil)
	db.EXPECT().ListEvents(mock.Anything).Return([]database.Event{*testDatabaseEvent()}, nil)
	outputService.EXPECT().ToLocalTime(mock.Anything, output).Return(time.Now())
	outputService.EXPECT().SendDailyReminder(&daemon.DailyReminder{
		Events: []daemon.Event{*testEvent()},
	}, output).Return(nil)
	db.EXPECT().UpdateOutput(mock.Anything).Return(nil, nil)

	go service.Start() //nolint:errcheck

	time.Sleep(time.Millisecond * 5)

	err := service.Stop()
	require.NoError(t, err)
}
