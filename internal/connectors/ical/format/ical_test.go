package format_test

import (
	"os"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"gorm.io/gorm"
)

func testTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2120-01-02T15:04:05+00:00")
	return t
}

func TestNewCalendar(t *testing.T) {
	events := []database.Event{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Time:           testTime(),
			RepeatUntil:    toP(testTime()),
			RepeatInterval: toP(time.Hour * 24),
			Message:        "Event 1",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Time:    testTime(),
			Message: "Event 2",
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			Time:           testTime(),
			RepeatUntil:    toP(testTime().Add(time.Hour * 100)),
			RepeatInterval: toP(time.Hour * 24),
			Message:        "Event 3",
		},
	}

	ical := format.NewCalendar("cal 1", events)
	icalShould, err := os.ReadFile("testdata/calendar1.ical")
	require.NoError(t, err)
	assert.Equal(t, string(icalShould), ical)
}

func TestMinutesToIcalRecurrenceRule(t *testing.T) {
	testCases := []struct {
		name         string
		interval     time.Duration
		occurences   uint64
		expectedRule string
	}{
		{
			name:       "4x 1 sec",
			interval:   time.Second,
			occurences: 4,
		},
		{
			name:         "4x 1 min",
			interval:     time.Minute,
			occurences:   4,
			expectedRule: "RRULE:FREQ=MINUTELY;COUNT=4",
		},
		{
			name:         "4x 2 min",
			interval:     time.Minute * 2,
			occurences:   4,
			expectedRule: "RRULE:FREQ=MINUTELY;INTERVAL=2;COUNT=4",
		},
		{
			name:         "4x 1 hour",
			interval:     time.Hour,
			occurences:   4,
			expectedRule: "RRULE:FREQ=HOURLY;COUNT=4",
		},
		{
			name:         "4x 2 hour",
			interval:     time.Hour * 2,
			occurences:   4,
			expectedRule: "RRULE:FREQ=HOURLY;INTERVAL=2;COUNT=4",
		},
		{
			name:         "4x 1 day",
			interval:     time.Hour * 24,
			occurences:   4,
			expectedRule: "RRULE:FREQ=DAILY;COUNT=4",
		},
		{
			name:         "4x 2 day",
			interval:     time.Hour * 48,
			occurences:   4,
			expectedRule: "RRULE:FREQ=DAILY;INTERVAL=2;COUNT=4",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.expectedRule,
				format.MinutesToIcalRecurrenceRule(tc.interval, tc.occurences),
			)
		})
	}
}

func TestEventsFromIcal(t *testing.T) {
	data, err := os.ReadFile("testdata/calendar1.ical")
	require.NoError(t, err)

	events, err := format.EventsFromIcal(string(data), &format.EventOpts{
		EventDelay: time.Duration(0),
	})
	require.NoError(t, err)

	require.Equal(t, 3, len(events))

	assert.Equal(t, testTime().UTC(), events[0].Time.UTC())
	assert.Equal(t, time.Minute*5, events[0].Duration)
	assert.Equal(t, "Event 1", events[0].Message)
	assert.Equal(t, "1", events[0].ExternalReference)
	assert.Equal(t, testTime().UTC(), events[1].Time.UTC())
	assert.Equal(t, time.Minute*5, events[1].Duration)
	assert.Equal(t, "Event 2", events[1].Message)
	assert.Equal(t, "2", events[1].ExternalReference)
	assert.Equal(t, testTime().UTC(), events[2].Time.UTC())
	assert.Equal(t, time.Minute*5, events[2].Duration)
	assert.Equal(t, "Event 3", events[2].Message)
	assert.Equal(t, "3", events[2].ExternalReference)
}

func TestEventsFromIcalWithAllDayEvent(t *testing.T) {
	data := `
BEGIN:VCALENDAR
VERSION:2.0
PRODID:cal 1
METHOD:PUBLISH
BEGIN:VEVENT
DTSTART:21200102Z
DTSTAMP:00010101T000000Z
UID:1
SUMMARY:Event 1
DESCRIPTION:Event 1
CLASS:PRIVATE
END:VEVENT
END:VCALENDAR
`

	events, err := format.EventsFromIcal(data, &format.EventOpts{
		EventDelay: time.Duration(0),
	})
	require.NoError(t, err)

	require.Equal(t, 1, len(events))

	assert.Equal(t, testTime().Round(time.Hour*24).Add(time.Hour*-24).UTC(), events[0].Time.UTC())
	assert.Equal(t, time.Hour*24, events[0].Duration)
	assert.Equal(t, "Event 1", events[0].Message)
	assert.Equal(t, "1", events[0].ExternalReference)
}

func TestEventsFromIcalWithNoTime(t *testing.T) {
	data := `
BEGIN:VCALENDAR
VERSION:2.0
PRODID:cal 1
METHOD:PUBLISH
BEGIN:VEVENT
DTSTAMP:00010101T000000Z
UID:1
SUMMARY:Event 1
DESCRIPTION:Event 1
CLASS:PRIVATE
END:VEVENT
END:VCALENDAR
`

	events, err := format.EventsFromIcal(data, &format.EventOpts{
		EventDelay: time.Duration(0),
	})
	require.NoError(t, err)

	require.Equal(t, 0, len(events))
}

func toP[T any](elem T) *T {
	return &elem
}
