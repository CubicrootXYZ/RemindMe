package daemon

func (service *service) sendOutEvents() error {
	events, err := service.database.GetEventsPending()
	if err != nil {
		return err
	}

	for _, event := range events {
		event := event
		for j := range event.Channel.Outputs {
			outputService, ok := service.config.OutputServices[event.Channel.Outputs[j].OutputType]

			if !ok {
				// TODO handle iCal here somehow
				service.logger.Errorf("missing output service for type: %s", event.Channel.Outputs[j].OutputType)
				continue
			}

			err = outputService.SendReminder(eventFromDatabase(&event), outputFromDatabase(&event.Channel.Outputs[j]))
			if err != nil {
				service.logger.Err(err)
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
				service.logger.Errorf("failed updating event after sending reminder: %w", err)
				continue
			}
		}
	}

	return nil
}
