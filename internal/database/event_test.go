package database_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func time2123() time.Time {
	t, _ := time.Parse(time.RFC3339, "2123-01-02T15:04:05-07:00")
	return t
}

func time2125() time.Time {
	t, _ := time.Parse(time.RFC3339, "2125-01-02T15:04:05-07:00")
	return t
}

func timeP2124() *time.Time {
	t, _ := time.Parse(time.RFC3339, "2124-01-02T15:04:05-07:00")
	return &t
}

func timeP2122() *time.Time {
	t, _ := time.Parse(time.RFC3339, "2122-01-02T15:04:05-07:00")
	return &t
}

func interval40Hours() *time.Duration {
	d := time.Hour * 40

	return &d
}

func TestEvent_NextEventTime(t *testing.T) {
	type testCase struct {
		Name         string
		Event        database.Event
		ExpectedTime time.Time
	}

	testCases := []testCase{
		{
			Name:         "Empty event",
			Event:        database.Event{},
			ExpectedTime: time.Time{},
		},
		{
			Name: "Empty repeat until",
			Event: database.Event{
				Time:           time2123(),
				RepeatInterval: interval40Hours(),
			},
			ExpectedTime: time.Time{},
		},
		{
			Name: "Empty repeat interval",
			Event: database.Event{
				Time:        time2123(),
				RepeatUntil: timeP2124(),
			},
			ExpectedTime: time.Time{},
		},
		{
			Name: "Repeat until reached",
			Event: database.Event{
				Time:           time2123(),
				RepeatUntil:    timeP2122(),
				RepeatInterval: interval40Hours(),
			},
			ExpectedTime: time.Time{},
		},
		{
			Name: "Success",
			Event: database.Event{
				Time:           time2123(),
				RepeatUntil:    timeP2124(),
				RepeatInterval: interval40Hours(),
			},
			ExpectedTime: time2123().Add(*interval40Hours()),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actualTime := testCase.Event.NextEventTime()
			assert.Equal(t, testCase.ExpectedTime, actualTime)
		})
	}
}

func testEvent() *database.Event {
	var channel *database.Channel
	err := gormDB.First(&channel).Error
	if err != nil {
		panic(err)
	}

	channel = testChannel()
	channel, err = service.NewChannel(channel)
	if err != nil {
		panic(err)
	}

	return &database.Event{
		Time:           time2123(),
		Duration:       *interval40Hours(),
		Message:        "test",
		Active:         true,
		RepeatInterval: interval40Hours(),
		RepeatUntil:    timeP2124(),
		ChannelID:      channel.ID,
	}
}

func TestService_NewEvent(t *testing.T) {
	eventBefore := testEvent()
	eventAfter, err := service.NewEvent(eventBefore)
	require.NoError(t, err)

	assert.Equal(t, eventBefore.Time, eventAfter.Time)
	assert.Equal(t, eventBefore.Duration, eventAfter.Duration)
	assert.Equal(t, eventBefore.Message, eventAfter.Message)
	assert.Equal(t, eventBefore.Active, eventAfter.Active)
	assert.Equal(t, eventBefore.RepeatInterval, eventAfter.RepeatInterval)
	assert.Equal(t, eventBefore.RepeatUntil, eventAfter.RepeatUntil)
}

func TestService_NewEvents(t *testing.T) {
	eventsBefore := []database.Event{
		*testEvent(),
		*testEvent(),
	}
	err := service.NewEvents(eventsBefore)
	require.NoError(t, err)

	require.Greater(t, eventsBefore[0].ID, uint(0))
	require.Greater(t, eventsBefore[1].ID, uint(0))
}

func TestService_NewEventWithoutChannel(t *testing.T) {
	_, err := service.NewEvent(&database.Event{})
	require.Error(t, err)
}

func TestService_GetEventsByChannel(t *testing.T) {
	eventBefore, err := service.NewEvent(testEvent())
	require.NoError(t, err)

	events, err := service.GetEventsByChannel(eventBefore.ChannelID)
	require.NoError(t, err)

	require.Less(t, 0, len(events))
	evtFound := false
	for _, eventAfter := range events {
		if eventAfter.ID == eventBefore.ID {
			evtFound = true
			assert.Equal(t, eventBefore.Duration, eventAfter.Duration)
			assert.Equal(t, eventBefore.Message, eventAfter.Message)
			assert.Equal(t, eventBefore.Active, eventAfter.Active)
			assert.Equal(t, eventBefore.RepeatInterval, eventAfter.RepeatInterval)
		}
	}

	assert.True(t, evtFound, "missing event not in response")
}

func TestService_ListEvents(t *testing.T) {
	eventBefore, err := service.NewEvent(testEvent())
	require.NoError(t, err)

	input := &database.Input{
		InputType: "test",
		InputID:   uint(rand.Int()), //nolint:gosec
	}
	err = service.AddInputToChannel(eventBefore.ChannelID, input)
	require.NoError(t, err)

	events, err := service.ListEvents(&database.ListEventsOpts{
		InputID: &input.ID,
	})
	require.NoError(t, err)

	require.Less(t, 0, len(events))
	evtFound := false
	for _, eventAfter := range events {
		if eventAfter.ID == eventBefore.ID {
			evtFound = true
			assert.Equal(t, eventBefore.Duration, eventAfter.Duration)
			assert.Equal(t, eventBefore.Message, eventAfter.Message)
			assert.Equal(t, eventBefore.Active, eventAfter.Active)
			assert.Equal(t, eventBefore.RepeatInterval, eventAfter.RepeatInterval)
		}
	}

	assert.True(t, evtFound, "missing event not in response")
}

func TestService_GetEventsByChannelWithEmptyResponse(t *testing.T) {
	events, err := service.GetEventsByChannel(123456)
	require.NoError(t, err)

	require.Equal(t, 0, len(events))
}

func TestService_GetEventsPending(t *testing.T) {
	eventBefore := testEvent()
	eventBefore.Time = time.Now().Add(-200 * time.Hour)

	eventBefore, err := service.NewEvent(eventBefore)
	require.NoError(t, err)

	events, err := service.GetEventsPending()
	require.NoError(t, err)

	require.Less(t, 0, len(events))
	evtFound := false
	for _, eventAfter := range events {
		if eventAfter.ID == eventBefore.ID {
			evtFound = true
			assert.Equal(t, eventBefore.Duration, eventAfter.Duration)
			assert.Equal(t, eventBefore.Message, eventAfter.Message)
			assert.Equal(t, eventBefore.Active, eventAfter.Active)
			assert.Equal(t, eventBefore.RepeatInterval, eventAfter.RepeatInterval)
		}
	}

	assert.True(t, evtFound, "missing event not in response")
}

func TestService_GetEventsPendingWithInactiveEvent(t *testing.T) {
	eventBefore := testEvent()
	eventBefore.Time = time.Now().Add(-200 * time.Hour)
	eventBefore.Active = false

	eventBefore, err := service.NewEvent(testEvent())
	require.NoError(t, err)

	events, err := service.GetEventsPending()
	require.NoError(t, err)

	evtFound := false
	for _, eventAfter := range events {
		if eventAfter.ID == eventBefore.ID {
			evtFound = true
		}
	}

	assert.False(t, evtFound, "event in response")
}

func TestService_UpdateEvent(t *testing.T) {
	eventBefore, err := service.NewEvent(testEvent())
	require.NoError(t, err)

	eventBefore.Time = time2125()
	eventBefore.Duration = time.Minute
	eventBefore.Message = "test 2"
	eventBefore.Active = false

	eventAfter, err := service.UpdateEvent(eventBefore)
	require.NoError(t, err)

	assert.Equal(t, eventBefore.Time, eventAfter.Time)
	assert.Equal(t, eventBefore.Duration, eventAfter.Duration)
	assert.Equal(t, eventBefore.Message, eventAfter.Message)
	assert.Equal(t, eventBefore.Active, eventAfter.Active)
	assert.Equal(t, eventBefore.RepeatInterval, eventAfter.RepeatInterval)
	assert.Equal(t, eventBefore.RepeatUntil, eventAfter.RepeatUntil)
}

func TestService_DeleteEvent(t *testing.T) {
	event, err := service.NewEvent(testEvent())
	require.NoError(t, err)

	err = service.DeleteEvent(event)
	require.NoError(t, err)

	events, err := service.ListEvents(&database.ListEventsOpts{
		IDs: []uint{event.ID},
	})
	require.NoError(t, err)
	assert.Len(t, events, 0)
}
