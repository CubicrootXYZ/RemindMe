package daemon_test

import (
	"errors"
	"testing"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon/mocks"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func testDaemon(ctrl *gomock.Controller, events, daily bool) (daemon.Service, *database.MockService, *mocks.MockOutputService) {
	db := database.NewMockService(ctrl)
	outputService := mocks.NewMockOutputService(ctrl)

	return daemon.New(&daemon.Config{
		OutputServices: map[string]daemon.OutputService{
			"test": outputService,
		},
		EventsInterval:        intervalFromBool(events),
		DailyReminderInterval: intervalFromBool(daily),
	}, db, gologger.New(gologger.LogLevelDebug, 0)), db, outputService
}

func intervalFromBool(b bool) time.Duration {
	if b {
		return time.Millisecond
	}

	return time.Hour
}

func TestMain(m *testing.M) {
	m.Run()
}

func time2123() time.Time {
	t, _ := time.Parse(time.RFC3339, "2123-01-02T15:04:05-07:00")
	return t
}

func timeP2122() *time.Time {
	t, _ := time.Parse(time.RFC3339, "2122-01-02T15:04:05-07:00")
	return &t
}

func interval40Hours() *time.Duration {
	d := time.Hour * 40

	return &d
}

func testDatabaseEvent() *database.Event {
	return &database.Event{
		Time:           time2123(),
		Message:        "test",
		RepeatInterval: interval40Hours(),
		RepeatUntil:    timeP2122(),
		Channel: database.Channel{
			Outputs: []database.Output{*testDatabaseOutput()},
		},
	}
}

func testDatabaseOutput() *database.Output {
	output := &database.Output{
		OutputType: "test",
		OutputID:   1,
	}

	output.ID = 12
	return output
}

func testEvent() *daemon.Event {
	return &daemon.Event{
		EventTime:      time2123(),
		Message:        "test",
		RepeatInterval: interval40Hours(),
		RepeatUntil:    timeP2122(),
	}
}

func testOutput() *daemon.Output {
	return &daemon.Output{
		ID:         12,
		OutputType: "test",
		OutputID:   1,
	}
}

func TestService_SendOutEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, db, outputService := testDaemon(ctrl, true, false)

	event := testDatabaseEvent()
	db.EXPECT().GetEventsPending().MinTimes(1).Return([]database.Event{*event}, nil)
	outputService.EXPECT().SendReminder(testEvent(), testOutput()).MinTimes(1).Return(nil)

	event.Active = false
	db.EXPECT().UpdateEvent(event).MinTimes(1).Return(nil, nil)

	go service.Start()               //nolint:errcheck
	time.Sleep(time.Millisecond * 5) // give time to execute

	err := service.Stop()
	require.NoError(t, err)
}

func TestService_SendOutEventsWithDatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, db, _ := testDaemon(ctrl, true, false)

	db.EXPECT().GetEventsPending().MinTimes(1).Return([]database.Event{*testDatabaseEvent()}, errors.New("test"))

	go service.Start()               //nolint:errcheck
	time.Sleep(time.Millisecond * 5) // give time to execute

	err := service.Stop()
	require.NoError(t, err)
}

func TestService_SendOutEventsWithOutputError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, db, outputService := testDaemon(ctrl, true, false)

	event := testDatabaseEvent()
	db.EXPECT().GetEventsPending().MinTimes(1).Return([]database.Event{*testDatabaseEvent()}, nil)
	outputService.EXPECT().SendReminder(testEvent(), testOutput()).MinTimes(1).Return(errors.New("test"))

	event.Active = false
	db.EXPECT().UpdateEvent(event).MinTimes(1).Return(nil, nil)

	go service.Start()               //nolint:errcheck
	time.Sleep(time.Millisecond * 5) // give time to execute

	err := service.Stop()
	require.NoError(t, err)
}
