package format_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageFromEvent(t *testing.T) {
	msg, msgF, err := format.MessageFromEvent(&daemon.Event{
		Message:   "my event",
		ID:        1,
		EventTime: refTime(),
	}, "")
	require.NoError(t, err)
	assert.Equal(
		t,
		`🔔 NEW EVENT:"my event"

ID: 1; Scheduled for 11:45 12.11.2014 (UTC) `,
		msg,
	)
	assert.Equal(
		t,
		`🔔 <b>New Event:</b>"my event"<br><br><i>ID: 1; </i><i>Scheduled for 11:45 12.11.2014 (UTC)</i> `,
		msgF,
	)
}
func TestMessageFromEventWithRecurring(t *testing.T) {
	dur := time.Hour
	msg, msgF, err := format.MessageFromEvent(&daemon.Event{
		Message:        "my event",
		ID:             1,
		EventTime:      refTime(),
		RepeatInterval: &dur,
	}, "")
	require.NoError(t, err)
	assert.Equal(
		t,
		`🔔 NEW EVENT:"my event"

ID: 1; Scheduled for 11:45 12.11.2014 (UTC) 🔁`,
		msg,
	)
	assert.Equal(
		t,
		`🔔 <b>New Event:</b>"my event"<br><br><i>ID: 1; </i><i>Scheduled for 11:45 12.11.2014 (UTC)</i> 🔁`,
		msgF,
	)
}

func TestInfoFromEvent(t *testing.T) {}
