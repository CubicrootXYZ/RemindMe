package calendar

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/stretchr/testify/assert"
)

var reminders *[]database.Reminder

func TestMain(m *testing.M) {
	r1 := database.Reminder{}
	r1.ID = 1
	r1.RemindTime = time.Now()
	r1.Message = "Hello World"
	r1.Active = true
	r1.RepeatInterval = 14
	r1.RepeatMax = 10

	reminderList := make([]database.Reminder, 0)
	reminderList = append(reminderList, r1)

	reminders = &reminderList

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCalendar_ICal(t *testing.T) {
	assert := assert.New(t)

	calendar := NewCalendar(reminders)
	ical := calendar.ICal()
	//log.Info(ical) // Used for debugging

	assert.Greaterf(len(ical), 20, "Calendar is <= 20 characters. That can not be!")
	for _, reminder := range *reminders {
		assert.Truef(strings.Contains(ical, reminder.Message), "Missing message \"%s\" in ical output", reminder.Message)
		assert.Truef(strings.Contains(ical, reminder.RemindTime.Format(formater.DateFormatICal)), "Missing date in ical output")
	}

	assert.Equalf(len(*reminders), strings.Count(ical, "BEGIN:VEVENT"), "Expected %d event begins but got %d", len(*reminders), strings.Count(ical, "BEGIN:VEVENT"))
	assert.Equalf(len(*reminders), strings.Count(ical, "END:VEVENT"), "Expected %d event ends but got %d", len(*reminders), strings.Count(ical, "END:VEVENT"))
}
