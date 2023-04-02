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
	t, _ := time.Parse(time.RFC3339, "2120-01-02T15:04:05+07:00")
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

func toP[T any](elem T) *T {
	return &elem
}
