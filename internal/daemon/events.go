package daemon

func (service *service) sendOutEvents() error {
	events, err := service.database.GetEventsPending()
	if err != nil {
		return err
	}

	for _, event := range events {
		eventSuccess := true
		for j := range event.Channel.Outputs {
			outputService, ok := service.config.OutputServices[event.Channel.Outputs[j].OutputType]

			if !ok {
				service.logger.Errorf("missing output service for type: %s", event.Channel.Outputs[j].OutputType)
				continue
			}

			err = outputService.SendReminder(eventFromDatabase(&event), outputFromDatabase(&event.Channel.Outputs[j])) //nolint:gosec // Reference stays in same routine.
			if err != nil {
				eventSuccess = false
				service.logger.Err(err)
				continue
			}
		}

		if !eventSuccess {
			// Try again later.
			//
			// This approach is only viable as there is only 1 output (matrix channel) possible currently. The iCal
			// output is no-op. Adding more outputs requires another solution here!
			continue
		}

		nextTime := event.NextEventTime()
		if !nextTime.IsZero() {
			event.Time = nextTime
		} else {
			event.Active = false
		}

		_, err = service.database.UpdateEvent(&event) //nolint:gosec // Reference stays in same routine.
		if err != nil {
			service.logger.Errorf("failed updating event after sending reminder: %w", err)
			continue
		}
	}

	return nil
}
