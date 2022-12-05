package formater

import (
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestReminderToMessage(t *testing.T) {
	reminder := &database.Reminder{
		Message:    "test reminder",
		RemindTime: refTime(),
		Channel: database.Channel{
			UserIdentifier: "@testuser:example.com",
		},
	}

	message, messageFormatted := ReminderToMessage(reminder)

	assert.Equal(t, "<a href=\"https://matrix.to/#/@testuser:example.com\">@testuser</a> a Reminder for you: <br>test reminder <br><i>(at 11:45 12.11.2014 (UTC))</i>", messageFormatted)
	assert.Equal(t, "testuser a reminder for you: test reminder  (at 11:45 12.11.2014 (UTC))", message)
}
