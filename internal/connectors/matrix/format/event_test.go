package format_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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
		`🔔 MY EVENT (#1)
11:45 (UTC) `,
		msg,
	)
	assert.Equal(
		t,
		`🔔 <b>my event</b> (#1)<br><i>11:45 (UTC)</i> `,
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
		`🔔 MY EVENT (#1)
11:45 (UTC) 🔁`,
		msg,
	)
	assert.Equal(
		t,
		`🔔 <b>my event</b> (#1)<br><i>11:45 (UTC)</i> 🔁`,
		msgF,
	)
}

func TestInfoFromEvent(t *testing.T) {
	msg, msgF := format.InfoFromEvent(&database.Event{
		Model: gorm.Model{
			ID: 1,
		},
		Message: "my event",
		Time:    refTime(),
	}, "")

	assert.Equal(
		t,
		`➡️ MY EVENT
at 11:45 12.11.2014 (UTC) (ID: 1) 
`,
		msg,
	)
	assert.Equal(
		t,
		`➡️ <b>my event</b><br>at 11:45 12.11.2014 (UTC) (ID: 1) <br>`,
		msgF,
	)
}

func TestInfoFromEventWithRecurring(t *testing.T) {
	dur := time.Hour
	msg, msgF := format.InfoFromEvent(&database.Event{
		Model: gorm.Model{
			ID: 1,
		},
		Message:        "my event",
		Time:           refTime(),
		RepeatInterval: &dur,
	}, "")

	assert.Equal(
		t,
		`➡️ MY EVENT
at 11:45 12.11.2014 (UTC) (ID: 1) 🔁 
`,
		msg,
	)
	assert.Equal(
		t,
		`➡️ <b>my event</b><br>at 11:45 12.11.2014 (UTC) (ID: 1) <i>🔁 </i><br>`,
		msgF,
	)
}

func TestInfoFromEvents(t *testing.T) {
	msg, msgF := format.InfoFromEvents([]database.Event{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Message: "my event 1",
			Time:    refTime(),
		},
		{
			Model: gorm.Model{
				ID: 1,
			},
			Message: "my event 2",
			Time:    refTime().Add(time.Minute * 2),
		},
		{
			Model: gorm.Model{
				ID: 1,
			},
			Message: "my event 3",
			Time:    refTime().Add(time.Minute * -2),
		},
	}, "")

	assert.Equal(
		t,
		"\nNOVEMBER\n➡️ MY EVENT 3\nat 11:43 12.11.2014 (UTC) (ID: 1) \n➡️ MY EVENT 1\nat 11:45 12.11.2014 (UTC) (ID: 1) \n➡️ MY EVENT 2\nat 11:47 12.11.2014 (UTC) (ID: 1) \n",
		msg,
	)
	assert.Equal(
		t,
		"<br><b>November</b><br>\n➡️ <b>my event 3</b><br>at 11:43 12.11.2014 (UTC) (ID: 1) <br>➡️ <b>my event 1</b><br>at 11:45 12.11.2014 (UTC) (ID: 1) <br>➡️ <b>my event 2</b><br>at 11:47 12.11.2014 (UTC) (ID: 1) <br>",
		msgF,
	)
}

func TestInfoFromEventsWithNoEvent(t *testing.T) {
	msg, msgF := format.InfoFromEvents(nil, "")

	assert.Equal(
		t,
		`no pending events found`,
		msg,
	)
	assert.Equal(
		t,
		`<i>no pending events found</i>`,
		msgF,
	)
}

func TestInfoFromDaemonEvent(t *testing.T) {
	testCases := []struct {
		name                 string
		event                *daemon.Event
		timeZone             string
		expectedMsg          string
		expectedFormattedMsg string
	}{
		{
			name: "nil event",
		},
		{
			name: "simple event",
			event: &daemon.Event{
				Message:   "my event",
				EventTime: refTime(),
			},
			expectedMsg:          "➡️ MY EVENT\nat 11:45 12.11.2014 (UTC) (ID: 0) \n",
			expectedFormattedMsg: "➡️ <b>my event</b><br>at 11:45 12.11.2014 (UTC) (ID: 0) <br>",
		},
		{
			name: "simple event with repeat interval",
			event: &daemon.Event{
				Message:        "my event",
				EventTime:      refTime(),
				RepeatInterval: toP(time.Hour),
			},
			expectedMsg:          "➡️ MY EVENT\nat 11:45 12.11.2014 (UTC) (ID: 0) 🔁 \n",
			expectedFormattedMsg: "➡️ <b>my event</b><br>at 11:45 12.11.2014 (UTC) (ID: 0) <i>🔁 </i><br>",
		},
		{
			name: "simple event with timezone",
			event: &daemon.Event{
				Message:   "my event",
				EventTime: refTime(),
			},
			timeZone:             "America/New_York",
			expectedMsg:          "➡️ MY EVENT\nat 06:45 12.11.2014 (EST) (ID: 0) \n",
			expectedFormattedMsg: "➡️ <b>my event</b><br>at 06:45 12.11.2014 (EST) (ID: 0) <br>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg, msgFormatted := format.InfoFromDaemonEvent(tc.event, tc.timeZone)
			assert.Equal(t, tc.expectedMsg, msg)
			assert.Equal(t, tc.expectedFormattedMsg, msgFormatted)
		})
	}
}

func TestInfoFromDaemonEvents(t *testing.T) {
	msg, formattedMsg := format.InfoFromDaemonEvents(
		[]daemon.Event{
			{
				Message:        "my event",
				EventTime:      refTime(),
				RepeatInterval: toP(time.Hour),
			},
			{
				Message:   "my event2",
				EventTime: refTime(),
			},
		}, "",
	)

	assert.Equal(t, "➡️ MY EVENT\nat 11:45 12.11.2014 (UTC) (ID: 0) 🔁 \n➡️ MY EVENT2\nat 11:45 12.11.2014 (UTC) (ID: 0) \n", msg)
	assert.Equal(t, "➡️ <b>my event</b><br>at 11:45 12.11.2014 (UTC) (ID: 0) <i>🔁 </i><br>➡️ <b>my event2</b><br>at 11:45 12.11.2014 (UTC) (ID: 0) <br>", formattedMsg)
}

func TestInfoFromDaemonEventsWithNoEvents(t *testing.T) {
	msg, formattedMsg := format.InfoFromDaemonEvents(
		nil, "",
	)

	assert.Equal(t, "no pending events found", msg)
	assert.Equal(t, "<i>no pending events found</i>", formattedMsg)
}

func toP[T any](elem T) *T {
	return &elem
}
