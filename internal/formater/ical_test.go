package formater

import (
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestFormater_ReminderToIcalEventOnSuccess(t *testing.T) {
	refTime := refTime()

	reminder := database.Reminder{
		RemindTime:     refTime,
		Message:        "my message",
		RepeatInterval: 30,
		RepeatMax:      10,
	}
	reminder.CreatedAt = refTime.AddDate(0, 0, -1)

	icalShould := "BEGIN:VEVENT\nDTSTART:20141112T114526Z\nDTEND:20141112T115026Z\nDTSTAMP:20141111T114526Z\nRRULE:FREQ=MINUTELY;INTERVAL=30;COUNT=10\nUID:0\nSUMMARY:my message\nDESCRIPTION:my message\nCLASS:PRIVATE\nEND:VEVENT\n"

	ical := ReminderToIcalEvent(&reminder)
	assert.Equal(t, icalShould, ical)
}

func TestFormater_ReminderToIcalEventOnFailure(t *testing.T) {
	ical := ReminderToIcalEvent(nil)
	assert.Equal(t, "", ical)
}

func TestFormater_MinutesToIcalReccurenceRule(t *testing.T) {
	testCases := make([]testIcalRecurrency, 0)
	testCases = append(testCases, testIcalRecurrency{
		minutes:     46,
		occurences:  49568,
		returnValue: "RRULE:FREQ=MINUTELY;INTERVAL=46;COUNT=49568",
	})
	testCases = append(testCases, testIcalRecurrency{
		minutes:     60 * 24 * 5,
		occurences:  2,
		returnValue: "RRULE:FREQ=DAILY;INTERVAL=5;COUNT=2",
	})
	testCases = append(testCases, testIcalRecurrency{
		minutes:     60 * 50,
		occurences:  0,
		returnValue: "RRULE:FREQ=HOURLY;INTERVAL=50",
	})

	for _, testCase := range testCases {
		rule := MinutesToIcalRecurrenceRule(testCase.minutes, testCase.occurences)

		assert.Equal(t, testCase.returnValue, rule)
	}
}

type testIcalRecurrency struct {
	minutes     uint64
	occurences  uint64
	returnValue string
}
