package daemon

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

func eventFromDatabase(event *database.Event) *Event {
	return &Event{
		ID:             event.ID,
		EventTime:      event.Time,
		Message:        event.Message,
		RepeatInterval: event.RepeatInterval,
		RepeatUntil:    event.RepeatUntil,
	}
}
