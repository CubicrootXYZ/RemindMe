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
	service.metricLastDailyReminderRun.WithLabelValues().
		Set(float64(time.Now().Unix()))

	eventsAfter := time.Now()
	eventsBefore := eventsAfter.Add(time.Hour*24 + time.Second)

	for _, channel := range channels {
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
			outputService, ok := service.config.OutputServices[output.OutputType]
			if !ok {
				service.logger.Error("unknown output type", "output.type", output.OutputType)
				continue
			}

			if !isDailyReminderTimeReached(channel, &output, outputService) {
				service.logger.
					Debug("no daily reminder send out",
						"channel.id", channel.ID,
						"output.id", output.ID,
						"output.output_type", output.OutputType,
						"reason", "reminder time not reached")

				continue
			}

			if isDailyReminderSentToday(&output, outputService) {
				service.logger.
					Debug("no daily reminder send out",
						"channel.id", channel.ID,
						"output.id", output.ID,
						"output.output_type", output.OutputType,
						"reason", "already done")

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

func isDailyReminderTimeReached(
	channel database.Channel,
	output *database.Output,
	outputService OutputService,
) bool {
	if channel.DailyReminder == nil {
		return false
	}

	now := outputService.ToLocalTime(time.Now(), outputFromDatabase(output))

	return (now.Hour()*60 + now.Minute()) >= int(*channel.DailyReminder)
}

func isDailyReminderSentToday(
	output *database.Output,
	outputService OutputService,
) bool {
	if output.LastDailyReminder == nil {
		return false
	}

	now := outputService.ToLocalTime(time.Now(), outputFromDatabase(output))

	return now.Day() == output.LastDailyReminder.Day() && now.Month() == output.LastDailyReminder.Month() && now.Year() == output.LastDailyReminder.Year()
}
