package format

import (
	"fmt"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
)

// MessageFromEvent creates a nicely formatted matrix message for the given event.
func MessageFromEvent(event *daemon.Event, timeZone string) (string, string, error) {
	f := Formater{}

	f.Text("🔔 ")
	f.Bold("New Event:")
	f.Text("\"")
	f.Text(event.Message)
	f.TextLine("\"")
	f.NewLine()

	f.Italic(fmt.Sprintf("ID: %d; ", event.ID))
	f.Italic("Scheduled for " + ToLocalTime(event.EventTime, timeZone))
	f.Text(" ")

	if event.RepeatInterval != nil {
		f.Text("🔁")
	}

	msg, msgFormatted := f.Build()
	return msg, msgFormatted, nil
}
