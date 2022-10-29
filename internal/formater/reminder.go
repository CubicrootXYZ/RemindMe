package formater

import (
	"fmt"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

func ReminderToMessage(reminder *database.Reminder) (message, messageHTML string) {
	message = fmt.Sprintf("%s a reminder for you: %s (at %s)", GetUsernameFromUserIdentifier(reminder.Channel.UserIdentifier), reminder.Message, ToLocalTime(reminder.RemindTime, reminder.Channel.TimeZone))
	messageHTML = fmt.Sprintf("%s a Reminder for you: <br>%s <br><i>(at %s)</i>", GetMatrixLinkForUser(reminder.Channel.UserIdentifier), reminder.Message, ToLocalTime(reminder.RemindTime, reminder.Channel.TimeZone))
	return
}
