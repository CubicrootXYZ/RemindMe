package format

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestHeaderFromEvent(t *testing.T) {
	testCases := []struct {
		name        string
		event       *database.Event
		loc         *time.Location
		expectedOut string
	}{
		{
			name: "today in UTC",
			event: &database.Event{
				Time: time.Now(),
			},
			loc:         time.UTC,
			expectedOut: "Today (" + time.Now().Format(DateFormatShort) + ")",
		},
		{
			name: "tomorrow in UTC",
			event: &database.Event{
				Time: time.Now().Add(time.Hour * 24),
			},
			loc:         time.UTC,
			expectedOut: "Tomorrow (" + time.Now().Add(time.Hour*24).Format(DateFormatShort) + ")",
		},
		{
			name: "this week in UTC",
			event: &database.Event{
				Time: time.Now().Add(time.Hour * 24 * 2),
			},
			loc:         time.UTC,
			expectedOut: "This Week",
		},
		{
			name: "next week in UTC",
			event: &database.Event{
				Time: time.Now().Add(time.Hour * 24 * 8),
			},
			loc:         time.UTC,
			expectedOut: "Next Week",
		},
		{
			name: "random month in UTC",
			event: &database.Event{
				Time: time.Now().Add(time.Hour * 24 * 32),
			},
			loc:         time.UTC,
			expectedOut: time.Now().Add(time.Hour * 24 * 32).Month().String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := headerFromEvent(tc.event, tc.loc)

			assert.Equal(t, tc.expectedOut, out)
		})
	}
}
