package formater

import (
	"fmt"
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

func ReminderToMessage(reminder *database.Reminder) (message, messageHTML string) {
	message = fmt.Sprintf("%s a reminder for you: %s %s (at %s)", GetUsernameFromUserIdentifier(reminder.Channel.UserIdentifier), reminder.Message, strings.Join(reminder.GetReminderIcons(), " "), ToLocalTime(reminder.RemindTime, reminder.Channel.TimeZone))
	messageHTML = fmt.Sprintf("%s a Reminder for you: <br>%s %s<br><i>(at %s)</i>", GetMatrixLinkForUser(reminder.Channel.UserIdentifier), reminder.Message, strings.Join(reminder.GetReminderIcons(), ""), ToLocalTime(reminder.RemindTime, reminder.Channel.TimeZone))
	return
}
