package reminderdaemon

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"gorm.io/gorm"
)

// CheckForDailyReminder checks which daily reminder messages needs to be send out and sends them.
func (d *Daemon) CheckForDailyReminder() error {
	log.Info("Checking daily reminders")
	channels, err := d.Database.GetChannelList()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	for _, channel := range channels {
		if channel.DailyReminder == nil {
			continue
		}

		tz := channel.Timezone()
		now := time.Now().In(tz)

		// We did not yet pass the daily reminder time
		if (now.Hour()*60 + now.Minute()) < int(*channel.DailyReminder) {
			continue
		}

		lastMessage, err := d.Database.GetLastMessageByType(database.MessageTypeDailyReminder, &channel)
		if err != nil && err != gorm.ErrRecordNotFound {
			continue
		}

		// Check if we already sent the reminder message
		if time.Since(lastMessage.CreatedAt) < 23*time.Hour+58*time.Minute {
			continue
		}

		dailyReminder, err := d.Database.GetDailyReminder(&channel)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		if len(*dailyReminder) == 0 {
			continue
		}

		log.Info(fmt.Sprintf("Sending out daily reminder to channel id %d", channel.ID))

		msg := &formater.Formater{}
		msg.Title("Your reminders for today")
		for _, reminder := range *dailyReminder {
			msg.BoldLine(reminder.Message)
			msg.Text("At ")
			msg.Text(formater.ToLocalTime(reminder.RemindTime, &channel))
			if reminder.Repeated != nil && reminder.RepeatMax > *reminder.Repeated {
				msg.ItalicLine(" (repeat every " + formater.ToNiceDuration(time.Minute*time.Duration(reminder.RepeatInterval)) + ")")
			} else {
				msg.NewLine()
			}
			msg.NewLine()
		}

		body, bodyFormatted := msg.Build()
		_, err = d.Messenger.SendFormattedMessage(body, bodyFormatted, &channel, database.MessageTypeDailyReminder, 0)
		if err != nil {
			log.Error(err.Error())
			continue
		}
	}

	return nil
}
