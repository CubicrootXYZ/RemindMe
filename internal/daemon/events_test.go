package daemon_test

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/require"
)

func testDaemon(t *testing.T, events, daily, cleanup bool) (daemon.Service, *database.MockService, *daemon.MockOutputService) {
	t.Helper()
	db := database.NewMockService(t)
	outputService := daemon.NewMockOutputService(t)

	return daemon.New(&daemon.Config{
		OutputServices: map[string]daemon.OutputService{
			"test": outputService,
		},
		EventsInterval:        intervalFromBool(events),
		DailyReminderInterval: intervalFromBool(daily),
		CleanupInterval:       intervalFromBool(cleanup),
	}, db, slog.Default()), db, outputService
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
	service, db, outputService := testDaemon(t, true, false, false)

	event := testDatabaseEvent()
	db.EXPECT().GetEventsPending().Return([]database.Event{*event}, nil)
	outputService.EXPECT().SendReminder(testEvent(), testOutput()).Return(nil)

	event.Active = false
	db.EXPECT().UpdateEvent(event).Return(nil, nil)

	go service.Start() //nolint:errcheck

	time.Sleep(time.Millisecond * 5) // give time to execute

	err := service.Stop()
	require.NoError(t, err)
}

func TestService_SendOutEventsWithDatabaseError(t *testing.T) {
	service, db, _ := testDaemon(t, true, false, false)

	db.EXPECT().GetEventsPending().Return([]database.Event{*testDatabaseEvent()}, errors.New("test"))

	go service.Start() //nolint:errcheck

	time.Sleep(time.Millisecond * 5) // give time to execute

	err := service.Stop()
	require.NoError(t, err)
}

func TestService_SendOutEventsWithOutputError(t *testing.T) {
	service, db, outputService := testDaemon(t, true, false, false)

	db.EXPECT().GetEventsPending().Return([]database.Event{*testDatabaseEvent()}, nil)
	outputService.EXPECT().SendReminder(testEvent(), testOutput()).Return(errors.New("test"))

	go service.Start() //nolint:errcheck

	time.Sleep(time.Millisecond * 5) // give time to execute

	err := service.Stop()
	require.NoError(t, err)
}
