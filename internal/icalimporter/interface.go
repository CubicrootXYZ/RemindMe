package icalimporter

import "errors"

var (
	// ErrMissingDtStart indicates that the event has no DTSTART entry
	ErrMissingDtStart = errors.New("missing DTSTART in event")
)

// IcalImporter imports events from ical sources and sets them as proper reminders
type IcalImporter interface {
	Run()
	Stop()
}
