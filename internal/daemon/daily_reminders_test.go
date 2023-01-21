package daemon_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, db, outputService := testDaemon(ctrl, false, true)

	channel := testDatabaseChannel()
	output := testOutput()

	db.EXPECT().GetChannels().MinTimes(1).Return([]database.Channel{channel}, nil)
	db.EXPECT().GetEventsByChannel(channel.ID).MinTimes(1).Return([]database.Event{*testDatabaseEvent()}, nil)
	outputService.EXPECT().SendDailyReminder(&daemon.DailyReminder{
		Events: []daemon.Event{*testEvent()},
	}, output).MinTimes(1).Return(nil)
	db.EXPECT().UpdateOutput(gomock.Any()).MinTimes(1).Return(nil, nil)

	go service.Start() //nolint:errcheck
	time.Sleep(time.Millisecond * 5)

	err := service.Stop()
	require.NoError(t, err)
}
