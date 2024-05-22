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

	service.logger.
		WithField("channels", len(channels)).
		Debugf("checking channels for daily reminder")

	for _, channel := range channels {
		if !isDailyReminderTimeReached(channel) {
			service.logger.
				WithField("channel", channel.ID).
				Debugf("no daily reminder send out - reminder time not reached")
			continue
		}

		events, err := service.database.GetEventsByChannel(channel.ID)
		if err != nil {
			service.logger.Err(err)
			continue
		}

		for _, output := range channel.Outputs {
			if isDailyReminderSentToday(&output) { //nolint:gosec // Reference stays in same routine.
				service.logger.
					WithField("channel", channel.ID).
					Debugf("no daily reminder send out - already done")
				continue
			}

			outputService, ok := service.config.OutputServices[output.OutputType]
			if !ok {
				service.logger.Errorf("missing output service for type: %s", output.OutputType)
				continue
			}

			service.logger.WithFields(map[string]any{
				"events":              len(events),
				"output_type":         output.OutputType,
				"output_id":           output.OutputID,
				"daily_reminder_time": channel.DailyReminder,
			}).Debugf("sending out daily reminder")

			err := outputService.SendDailyReminder(dailyReminderFromDatabase(events), outputFromDatabase(&output)) //nolint:gosec // Reference stays in same routine.
			if err != nil {
				service.logger.Err(err)
				continue
			}

			now := time.Now().UTC()
			output.LastDailyReminder = &now
			_, err = service.database.UpdateOutput(&output) //nolint:gosec // Reference stays in same routine.
			if err != nil {
				service.logger.Err(err)
			}
		}
	}

	return nil
}

func isDailyReminderTimeReached(channel database.Channel) bool {
	if channel.DailyReminder == nil {
		return false
	}

	now := time.Now().In(time.UTC)
	return (now.Hour()*60 + now.Minute()) >= int(*channel.DailyReminder)
}

func isDailyReminderSentToday(output *database.Output) bool {
	if output.LastDailyReminder == nil {
		return false
	}

	now := time.Now().In(time.UTC)
	return now.Day() == output.LastDailyReminder.Day() && now.Month() == output.LastDailyReminder.Month() && now.Year() == output.LastDailyReminder.Year()
}
