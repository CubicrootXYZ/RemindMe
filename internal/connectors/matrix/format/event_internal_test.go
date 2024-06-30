package format

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderFromEvent(t *testing.T) {
	baseTime, err := time.Parse(time.RFC3339, "2024-06-24T10:04:05Z")
	require.NoError(t, err)

	today := baseTime
	tomorrow := baseTime.Add(time.Hour * 24)
	thisWeek := baseTime.Add(time.Hour * 24 * 2)
	nextWeek := baseTime.Add(time.Hour * 24 * 7)
	nextMonth := baseTime.Add(time.Hour * 24 * 14)

	testCases := []struct {
		name        string
		event       *database.Event
		loc         *time.Location
		expectedOut string
	}{
		{
			name: "today in UTC",
			event: &database.Event{
				Time: today,
			},
			loc:         time.UTC,
			expectedOut: "Today (" + today.Format(DateFormatShort) + ")",
		},
		{
			name: "tomorrow in UTC",
			event: &database.Event{
				Time: tomorrow,
			},
			loc:         time.UTC,
			expectedOut: "Tomorrow (" + tomorrow.Format(DateFormatShort) + ")",
		},
		{
			name: "this week in UTC",
			event: &database.Event{
				Time: thisWeek,
			},
			loc:         time.UTC,
			expectedOut: "This Week",
		},
		{
			name: "next week in UTC",
			event: &database.Event{
				Time: nextWeek,
			},
			loc:         time.UTC,
			expectedOut: "Next Week",
		},
		{
			name: "random month in UTC",
			event: &database.Event{
				Time: nextMonth,
			},
			loc:         time.UTC,
			expectedOut: nextMonth.Month().String() + " 2024",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := headerFromEvent(tc.event, tc.loc, baseTime)

			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestHeaderFromEventWithTimeZone(t *testing.T) {
	tz, err := time.LoadLocation("Europe/Berlin")
	require.NoError(t, err)

	baseTime, err := time.Parse(time.RFC3339, "2024-06-30T13:39:05+02:00")
	require.NoError(t, err)

	eventTime, err := time.Parse(time.RFC3339, "2024-07-01T20:00:00+02:00")
	require.NoError(t, err)

	out := headerFromEvent(&database.Event{
		Time: eventTime,
	}, tz, baseTime)

	assert.Equal(t, "Tomorrow (Mon, 01 Jul)", out)
}
