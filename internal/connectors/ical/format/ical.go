package format

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	ical "github.com/arran4/golang-ical"
	"github.com/teambition/rrule-go"
)

// List of commonly used errors in this package.
var (
	ErrMissingDtStart     = errors.New("missing property DTSTART")
	ErrCanNotGetStartTime = errors.New("can not get event start time")
	ErrCanNotGetEndTime   = errors.New("can not get event end time")
)

var dateFormatICal = "20060102T150405Z07:00"

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
		ical.WriteString(event.Time.UTC().Format(dateFormatICal))
		ical.WriteString("\nDTEND:")
		ical.WriteString(endTime.UTC().Format(dateFormatICal))
		ical.WriteString("\nDTSTAMP:")
		ical.WriteString(event.CreatedAt.UTC().Format(dateFormatICal))
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

// EventOpts are options to apply to events.
type EventOpts struct {
	EventDelay      time.Duration // Delay to add to events.
	DefaultDuration time.Duration

	// If set the events will be pre-populated with this data.
	ChannelID uint
	InputID   uint
}

// EventsFromIcal extracts events from the iCal input.
func EventsFromIcal(input string, opts *EventOpts) ([]database.Event, error) {
	// TODO currently does not set recurring attributes.
	calendar, err := ical.ParseCalendar(strings.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("can not parse iCal calendar: %w", err)
	}

	events := make([]database.Event, 0)
	for _, event := range calendar.Events() {
		idProp := event.GetProperty(ical.ComponentPropertyUniqueId)
		if idProp == nil {
			continue
		}
		id := idProp.Value
		if len(id) <= 0 {
			continue
		}

		startTime, err := getStartTimeFromEvent(event)
		if err != nil {
			continue
		}
		duration, err := getDurationFromEvent(event)
		if err != nil {
			duration = opts.DefaultDuration
		}

		startTime = startTime.Add(opts.EventDelay)

		if time.Until(startTime) < 0 {
			// Ignore past events
			continue
		}

		name := getMessageFromEvent(event)
		if name == "" {
			continue
		}

		events = append(events, database.Event{
			Time:              startTime,
			Duration:          duration,
			Message:           getMessageFromEvent(event),
			Active:            true,
			ChannelID:         opts.ChannelID,
			InputID:           &opts.InputID,
			ExternalReference: id,
		})
	}

	return events, nil
}

func getStartTimeFromEvent(event *ical.VEvent) (time.Time, error) {
	startTime, err := event.GetStartAt()
	if err != nil {
		startTime, err = event.GetAllDayStartAt()
		if err != nil {
			return time.Now(), fmt.Errorf("%w: %w", ErrCanNotGetStartTime, err)
		}
	}

	rruleString := event.GetProperty(ical.ComponentPropertyRrule)
	if rruleString != nil {
		// RRULE needs the DTSTART too
		dtStart := event.GetProperty(ical.ComponentPropertyDtStart)
		if dtStart == nil {
			return time.Now(), ErrMissingDtStart
		}

		rruleObj, err := rrule.StrToRRule("DTSTART:" + dtStart.Value + "\n" + rruleString.Value)
		if err != nil {
			return time.Now(), fmt.Errorf("failed to parse RRULE: %w", err)
		}

		refTime := time.Now()
		if refTime.Sub(startTime) < 0 {
			refTime = startTime.Add(time.Second * -1)
		}
		return rruleObj.After(refTime, false), nil
	}

	return startTime, nil
}

func getDurationFromEvent(event *ical.VEvent) (time.Duration, error) {
	// Check if is usual event
	startTime, err := event.GetStartAt()
	if err == nil {
		endTime, err := event.GetEndAt()
		if err == nil {
			return endTime.Sub(startTime), nil
		}
	}

	// Else it might be an all day event?
	startTime, err = event.GetAllDayStartAt()
	if err == nil {
		endTime, err := event.GetAllDayEndAt()
		if err == nil {
			return endTime.Sub(startTime), nil
		}
		return 24 * time.Hour, nil
	}

	return time.Duration(0), ErrCanNotGetEndTime
}

func getMessageFromEvent(event *ical.VEvent) string {
	for _, property := range []ical.ComponentProperty{ical.ComponentPropertyDescription, ical.ComponentPropertySummary} {
		prop := event.GetProperty(property)
		if prop == nil {
			continue
		}

		return prop.Value
	}

	return ""
}
