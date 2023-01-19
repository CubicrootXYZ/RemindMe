package daemon

func (service *service) sendOutEvents() error {
	events, err := service.Database.GetEventsPending()
	if err != nil {
		return err
	}

	for i := range events {
		for j := range events[i].Channel.Outputs {
			outputService, ok := service.Config.OutputServices[events[i].Channel.Outputs[j].OutputType]

			if !ok {
				service.Logger.Errorf("missing output service for type: %s", events[i].Channel.Outputs[j].OutputType)
				continue
			}

			err = outputService.SendReminder(eventFromDatabase(&events[i]), outputFromDatabase(&events[i].Channel.Outputs[j]))
			if err != nil {
				service.Logger.Err(err)
				continue
			}
		}
	}

	return nil
}
