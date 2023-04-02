package format

import (
	"strconv"
	"strings"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

var dateFormatICal = "20060102T150405Z"

// NewCalendar creates an iCal formatted calendar with the given data.
func NewCalendar(calendarID string, events []database.Event) string {
	ical := strings.Builder{}
	ical.WriteString("BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:")
	ical.WriteString(calendarID)
	ical.WriteString("\nMETHOD:PUBLISH\n")

	for _, event := range events {
		// TODO make configurable
		endTime := event.Time.Add(time.Minute * 5)

		ical.WriteString("BEGIN:VEVENT\nDTSTART:")
		ical.WriteString(event.Time.Format(dateFormatICal))
		ical.WriteString("\nDTEND:")
		ical.WriteString(endTime.Format(dateFormatICal))
		ical.WriteString("\nDTSTAMP:")
		ical.WriteString(event.CreatedAt.Format(dateFormatICal))
		if event.RepeatInterval != nil {
			if event.RepeatUntil != nil && time.Until(*event.RepeatUntil) > 0 ||
				event.RepeatUntil == nil {
				ical.WriteString("\n")
				ical.WriteString(MinutesToIcalRecurrenceRule(
					*event.RepeatInterval,
					occurencesFromStartAndEnd(event.Time, *event.RepeatInterval, event.RepeatUntil),
				))
			}
		}
		ical.WriteString("\nUID:")
		ical.WriteString(strconv.FormatUint(uint64(event.ID), 10))
		ical.WriteString("\nSUMMARY:")
		ical.WriteString(event.Message)
		ical.WriteString("\nDESCRIPTION:")
		ical.WriteString(event.Message)
		ical.WriteString("\nCLASS:PRIVATE\nEND:VEVENT\n")
	}

	ical.WriteString("END:VCALENDAR\n")

	return ical.String()
}

// MinutesToIcalRecurrenceRule transfers the given minutes into an iCal recurrence rule
// https://icalendar.org/iCalendar-RFC-5545/3-8-5-3-recurrence-rule.html
func MinutesToIcalRecurrenceRule(interval time.Duration, occurences uint64) string {
	if interval.Minutes() < 1 {
		return ""
	}

	minutes := uint64(interval.Minutes())

	rule := strings.Builder{}
	rule.WriteString("RRULE:")

	switch {
	case minutes%(60*24) == 0:
		days := minutes / (60 * 24)
		rule.WriteString("FREQ=DAILY")
		if days > 1 {
			rule.WriteString(";INTERVAL=")
			rule.WriteString(strconv.FormatUint(days, 10))
		}
	case minutes%(60) == 0:
		hours := minutes / (60)
		rule.WriteString("FREQ=HOURLY")
		if hours > 1 {
			rule.WriteString(";INTERVAL=")
			rule.WriteString(strconv.FormatUint(hours, 10))
		}
	default:
		rule.WriteString("FREQ=MINUTELY")
		if minutes > 1 {
			rule.WriteString(";INTERVAL=")
			rule.WriteString(strconv.FormatUint(minutes, 10))
		}
	}

	if occurences > 0 {
		rule.WriteString(";COUNT=")
		rule.WriteString(strconv.FormatUint(occurences, 10))
	}

	return rule.String()
}

func occurencesFromStartAndEnd(start time.Time, interval time.Duration, end *time.Time) uint64 {
	if end == nil {
		return 0
	}

	occurences := uint64(0)
	start = start.Add(interval)
	for end.Sub(start) > 0 {
		start = start.Add(interval)
		occurences++
	}

	return occurences
}
