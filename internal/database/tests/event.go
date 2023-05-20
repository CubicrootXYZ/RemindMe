package tests

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"gorm.io/gorm"
)

func refTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	return t
}

type EventOpts func(*database.Event)

func TestEvent(opts ...EventOpts) database.Event {
	event := database.Event{
		Model: gorm.Model{
			ID: 2824,
		},
		Duration: time.Hour,
		Message:  "test event",
		Time:     refTime(),
	}

	for _, opt := range opts {
		opt(&event)
	}

	return event
}
