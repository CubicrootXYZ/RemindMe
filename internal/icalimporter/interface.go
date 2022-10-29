package icalimporter

// IcalImporter imports events from ical sources and sets them as proper reminders
type IcalImporter interface {
	Run()
	Stop()
}
