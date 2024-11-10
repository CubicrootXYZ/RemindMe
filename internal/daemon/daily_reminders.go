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
		With("channels", len(channels)).
		Debug("checking channels for daily reminder")

	eventsAfter := time.Now()
	eventsBefore := eventsAfter.Add(time.Hour*24 + time.Second)

	for _, channel := range channels {
		if !isDailyReminderTimeReached(channel) {
			service.logger.
				Debug("no daily reminder send out",
					"channel.id", channel.ID,
					"reason", "reminder time not reached")
			continue
		}

		events, err := service.database.ListEvents(&database.ListEventsOpts{
			ChannelID:    &channel.ID,
			EventsAfter:  &eventsAfter,
			EventsBefore: &eventsBefore,
		})
		if err != nil {
			service.logger.Error("failed to list events", "error", err)
			continue
		}

		for _, output := range channel.Outputs {
			if isDailyReminderSentToday(&output) {
				service.logger.
					Debug("no daily reminder send out",
						"channel.id", channel.ID,
						"reason", "already done")
				continue
			}

			outputService, ok := service.config.OutputServices[output.OutputType]
			if !ok {
				service.logger.Error("unknown output type", "output.type", output.OutputType)
				continue
			}

			service.logger.With(
				"events", len(events),
				"output.type", output.OutputType,
				"output.id", output.OutputID,
				"daily_reminder_time", channel.DailyReminder,
			).Debug("sending out daily reminder")

			err := outputService.SendDailyReminder(dailyReminderFromDatabase(events), outputFromDatabase(&output))
			if err != nil {
				service.logger.Error("failed to send out daily reminder", "error", err)
				continue
			}

			now := time.Now().UTC()
			output.LastDailyReminder = &now
			_, err = service.database.UpdateOutput(&output)
			if err != nil {
				service.logger.Error("failed to update output", "error", err)
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
