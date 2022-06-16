package calendar

import (
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
)

// Calendar for creating ical files
type Calendar struct {
	events []database.Reminder
}

// NewCalendar creates a new Calendar struct
func NewCalendar(events *[]database.Reminder) *Calendar {
	return &Calendar{
		events: *events,
	}
}

// ICal returns the calendars events in ical format
func (calendar *Calendar) ICal() string {
	ical := strings.Builder{}
	ical.WriteString("BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:RemindMe\nMETHOD:PUBLISH\n")
	if calendar.events != nil {
		for i := range calendar.events {
			ical.WriteString(formater.ReminderToIcalEvent(&calendar.events[i]))
		}
	}
	ical.WriteString("END:VCALENDAR\n")

	return ical.String()
}
