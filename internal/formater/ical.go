package formater

import (
	"strconv"
	"strings"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

// MinutesToIcalRecurrenceRule transfers the given minutes into an iCal recurrence rule
// https://icalendar.org/iCalendar-RFC-5545/3-8-5-3-recurrence-rule.html
func MinutesToIcalRecurrenceRule(minutes uint64, occurences uint64) string {
	if minutes == 0 {
		return ""
	}

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

// ReminderToIcalEvent formats a reminder into an iCal event
func ReminderToIcalEvent(reminder *database.Reminder) string {
	if reminder == nil {
		return ""
	}
	ical := strings.Builder{}
	endTime := reminder.RemindTime.Add(5 * time.Minute)
	ical.WriteString("BEGIN:VEVENT\nDTSTART:")
	ical.WriteString(reminder.RemindTime.Format(DateFormatICal))
	ical.WriteString("\nDTEND:")
	ical.WriteString(endTime.Format(DateFormatICal))
	ical.WriteString("\nDTSTAMP:")
	ical.WriteString(reminder.CreatedAt.Format(DateFormatICal))
	if reminder.RepeatInterval > 0 && reminder.RepeatMax > 0 {
		ical.WriteString("\n")
		ical.WriteString(MinutesToIcalRecurrenceRule(reminder.RepeatInterval, reminder.RepeatMax))
	}
	ical.WriteString("\nUID:")
	ical.WriteString(strconv.FormatUint(uint64(reminder.ID), 10))
	ical.WriteString("\nSUMMARY:")
	ical.WriteString(reminder.Message)
	ical.WriteString("\nDESCRIPTION:")
	ical.WriteString(reminder.Message)
	ical.WriteString("\nCLASS:PRIVATE\nEND:VEVENT\n")

	return ical.String()
}
