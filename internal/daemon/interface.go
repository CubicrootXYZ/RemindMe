package daemon

import (
	"time"
)

// Service defines a daemon service interface.
// The daemon is responsible for sending out reminders via the defined outputs.
type Service interface {
	Start() error
	Stop() error
}

// Event holds information about a reminder.
type Event struct {
	ID             uint
	EventTime      time.Time
	Message        string
	RepeatInterval *time.Duration
	RepeatUntil    *time.Time
}

// DailyReminder holds information about a daily reminder.
type DailyReminder struct {
	Events []Event
}

// Output holds information about the output.
type Output struct {
	ID         uint
	OutputType string
	OutputID   uint
}
