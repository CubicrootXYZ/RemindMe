package daemon

import (
	"time"
)

// Service defines a daemon service interface.
// The daemon is responsible for sending out reminders via the defined outputs.
type Service interface {
}

// OutputService defines an interface for services handling outputs.
type OutputService interface {
	SendReminder(*Event, *Output) error
	SendDailyReminder(*DailyReminder, *Output) error
}

// Event holds information about a reminder.
type Event struct {
	EventTime      time.Time
	Message        string
	RepeatInterval time.Duration
	RepeatCountMax uint64
	RepeatCount    uint64
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
