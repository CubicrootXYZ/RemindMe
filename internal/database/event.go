package database

import "time"

// NextEventTime returns the next time the event will happen.
// If the RepeatUntil date is reached or the event is not recurring the zero time is returned.
func (event *Event) NextEventTime() time.Time {
	if event.RepeatInterval == nil || event.RepeatUntil == nil || time.Until(*event.RepeatUntil) < 0 {
		return time.Time{}
	}

	nextTime := event.Time.Add(*event.RepeatInterval)
	for time.Until(nextTime) < 0 {
		nextTime.Add(*event.RepeatInterval)
	}

	if nextTime.After(*event.RepeatUntil) {
		return time.Time{}
	}

	return nextTime
}

func (service *service) NewEvent(event *Event) (*Event, error) {
	err := service.db.Create(event).Error

	return event, err
}

func (service *service) GetEventsByChannel(channelID uint) ([]Event, error) {
	var events []Event

	err := service.db.Preload("Channel").Preload("Input").Find(&events, "events.channel_id = ?", channelID).Error
	return events, err
}

func (service *service) GetEventsPending() ([]Event, error) {
	var events []Event

	err := service.db.Preload("Channel").Preload("Input").Find(&events, "events.active = ? AND events.time <= ?", true, time.Now()).Error
	return events, err
}

func (service *service) UpdateEvent(event *Event) (*Event, error) {
	err := service.db.Save(event).Error

	return event, err
}
