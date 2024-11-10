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
				service.logger.Error("unknown output type", "output.type", event.Channel.Outputs[j].OutputType)
				continue
			}

			err = outputService.SendReminder(eventFromDatabase(&event), outputFromDatabase(&event.Channel.Outputs[j]))
			if err != nil {
				eventSuccess = false
				service.logger.Error("failed to send reminder to output",
					"error", err,
					"output.id", event.Channel.Outputs[j].OutputID,
					"output.type", event.Channel.Outputs[j].OutputType)
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

		_, err = service.database.UpdateEvent(&event)
		if err != nil {
			service.logger.Error("failed updating event after sending reminder", "error", err)
			continue
		}
	}

	return nil
}
