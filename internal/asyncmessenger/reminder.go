package asyncmessenger

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"gorm.io/gorm"
)

// TODO move this away from here
// Reminder holds all information about a reminder
type Reminder struct {
	gorm.Model
	RemindTime     time.Time
	Message        string
	Active         bool
	RepeatInterval uint64
	RepeatMax      uint64
	Repeated       *uint64
	ChannelID      uint
	Channel        Channel
}

func (reminder *Reminder) getRemindMessage() (message, messageFormatted string) {
	message = fmt.Sprintf("%s a reminder for you: %s (at %s)", "USER", reminder.Message, formater.ToLocalTime(reminder.RemindTime, reminder.Channel.TimeZone))
	messageFormatted = fmt.Sprintf("%s a Reminder for you: <br>%s <br><i>(at %s)</i>", formater.GetMatrixLinkForUser(reminder.Channel.UserIdentifier), reminder.Message, formater.ToLocalTime(reminder.RemindTime, reminder.Channel.TimeZone))
	return message, messageFormatted
}
