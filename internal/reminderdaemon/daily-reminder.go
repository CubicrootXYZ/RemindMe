package reminderdaemon

import (
	"fmt"
	"strings"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/random"
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

	for i := range channels {
		if channels[i].DailyReminder == nil {
			continue
		}

		tz := channels[i].Timezone()
		now := time.Now().In(tz)

		// We did not yet pass the daily reminder time
		if (now.Hour()*60 + now.Minute()) < int(*channels[i].DailyReminder) {
			continue
		}

		lastMessage, err := d.Database.GetLastMessageByType(database.MessageTypeDailyReminder, &channels[i])
		if err != nil && err != gorm.ErrRecordNotFound {
			continue
		}

		// Check if we already sent the reminder message

		if lastMessage.CreatedAt.In(tz).Day() == now.Day() {
			continue
		}

		dailyReminder, err := d.Database.GetDailyReminder(&channels[i])
		if err != nil {
			log.Error(err.Error())
			continue
		}

		log.Info(fmt.Sprintf("Sending out daily reminder to channel id %d", channels[i].ID))

		msg := &formater.Formater{}
		if len(*dailyReminder) > 0 {
			msg.Title("Your reminders for today")
		} else {
			msg.Text("Nothing to do today 🥳. " + random.MotivationalSentence())
		}

		for _, reminder := range *dailyReminder {
			msg.BoldLine(reminder.Message)
			msg.Text("At ")
			msg.Text(formater.ToLocalTime(reminder.RemindTime, channels[i].TimeZone))
			if reminder.RepeatInterval > 0 && (reminder.Repeated == nil || reminder.RepeatMax > *reminder.Repeated) {
				msg.Italic(" (repeat every " + formater.ToNiceDuration(time.Minute*time.Duration(reminder.RepeatInterval)) + ")")
			}
			msg.Text(" " + strings.Join(reminder.GetReminderIcons(), " "))
			msg.NewLine()
		}

		body, bodyFormatted := msg.Build()

		go func(body, bodyFormatted string, channel database.Channel) {
			resp, err := d.Messenger.SendMessage(asyncmessenger.HTMLMessage(
				body,
				bodyFormatted,
				channel.ChannelIdentifier,
			))
			if err != nil {
				log.Error(err.Error())
				return
			}
			_, err = d.Database.AddMessage(
				&database.Message{
					Body:               body,
					BodyHTML:           bodyFormatted,
					Type:               database.MessageTypeDailyReminder,
					ChannelID:          channel.ID,
					Timestamp:          resp.Timestamp,
					ExternalIdentifier: resp.ExternalIdentifier,
				},
			)
			if err != nil {
				log.Error(err.Error())
			}
		}(body, bodyFormatted, channels[i])
	}

	return nil
}
