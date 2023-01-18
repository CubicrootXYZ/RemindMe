package database_test

import (
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

func TestService_GetEventsByChannelWithEmptyResponse(t *testing.T) {
	eventBefore, err := service.NewEvent(testEvent())
	require.NoError(t, err)

	events, err := service.GetEventsByChannel(eventBefore.ChannelID)
	require.NoError(t, err)

	require.Equal(t, len(events), 0)
}

func TestService_GetEventsPending(t *testing.T) {
	eventBefore := testEvent()
	eventBefore.Time = time.Now().Add(-200 * time.Hour)

	eventBefore, err := service.NewEvent(testEvent())
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
