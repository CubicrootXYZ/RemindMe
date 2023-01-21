package daemon

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

func (service *service) sendOutDailyReminders() error {
	channels, err := service.database.GetChannels()
	if err != nil {
		return err
	}

	for _, channel := range channels {
		channel := channel

		if !isDailyReminderTimeReached(&channel) {
			continue
		}

		events, err := service.database.GetEventsByChannel(channel.ID)
		if err != nil {
			service.logger.Err(err)
			continue
		}

		for _, output := range channel.Outputs {
			output := output

			if isDailyReminderSentToday(&output, &channel) {
				continue
			}

			outputService, ok := service.config.OutputServices[output.OutputType]
			if !ok {
				service.logger.Errorf("missing output service for type: %s", output.OutputType)
				continue
			}

			err := outputService.SendDailyReminder(dailyReminderFromDatabase(events), outputFromDatabase(&output))
			if err != nil {
				service.logger.Err(err)
				continue
			}

			now := time.Now().UTC()
			output.LastDailyReminder = &now
			_, err = service.database.UpdateOutput(&output)
			if err != nil {
				service.logger.Err(err)
			}
		}
	}

	return nil
}

func isDailyReminderTimeReached(channel *database.Channel) bool {
	if channel.DailyReminder == nil {
		return false
	}

	now := time.Now().In(channel.Timezone())
	return (now.Hour()*60 + now.Minute()) >= int(*channel.DailyReminder)
}

func isDailyReminderSentToday(output *database.Output, channel *database.Channel) bool {
	if output.LastDailyReminder == nil {
		return false
	}

	now := time.Now().In(channel.Timezone())
	return now.Day() == output.LastDailyReminder.Day() && now.Month() == output.LastDailyReminder.Month() && now.Year() == output.LastDailyReminder.Year()
}
