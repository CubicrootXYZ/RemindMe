package format

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
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

// InfoFromEvent translates a database event into a nice human readable format.
func InfoFromEvent(event *database.Event, timeZone string) (string, string) {
	f := Formater{}
	f.Text("➡️ ")
	f.BoldLine(event.Message)
	f.Text("at ")
	f.Text(ToLocalTime(event.Time, timeZone))
	f.Text(" (ID: ")
	f.Text(strconv.Itoa(int(event.ID)))
	f.Text(") ")
	if event.RepeatInterval != nil {
		f.Italic("🔁 ")
	}
	if event.ExternalReference != "" {
		f.Italic("🌐 ")
	}
	f.NewLine()

	return f.Build()
}

// InfoFromEvents translates multiple database events into a nice human readable format.
func InfoFromEvents(events []database.Event, timeZone string) (string, string) {
	var str, strFormatted strings.Builder
	for i := range events {
		msg, msgF := InfoFromEvent(&events[i], timeZone)
		str.WriteString(msg)
		strFormatted.WriteString(msgF)
	}

	return str.String(), strFormatted.String()
}