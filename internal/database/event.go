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
		nextTime = nextTime.Add(*event.RepeatInterval)
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

func (service *service) NewEvents(events []Event) error {
	return service.db.Create(events).Error
}

// ListEventsOpts holds options for listing events.
type ListEventsOpts struct {
	IDs       []uint
	InputID   *uint
	ChannelID *uint

	IncludeInactive bool

	EventsBefore *time.Time
	EventsAfter  *time.Time
}

func (service *service) ListEvents(opts *ListEventsOpts) ([]Event, error) {
	query := service.db

	if opts.InputID != nil {
		query = query.Where("events.input_id = ?", *opts.InputID)
	}
	if opts.ChannelID != nil {
		query = query.Where("events.channel_id = ?", *opts.ChannelID)
	}
	if opts.IDs != nil {
		query = query.Where("events.id IN ?", opts.IDs)
	}
	if !opts.IncludeInactive {
		query = query.Where("events.active = ?", true)
	}
	if opts.EventsBefore != nil {
		query = query.Where("events.time <= ?", *opts.EventsBefore)
	}
	if opts.EventsAfter != nil {
		query = query.Where("events.time >= ?", *opts.EventsAfter)
	}

	query = query.Order("time DESC")

	var events []Event

	err := query.Preload("Channel").Preload("Input").Find(&events).Error
	return events, err
}

// TODO replace with ListEvents
func (service *service) GetEventsByChannel(channelID uint) ([]Event, error) {
	var events []Event

	err := service.db.Preload("Channel").Preload("Input").Find(&events, "events.channel_id = ?", channelID).Error
	return events, err
}

// TODO replace with ListEvents
func (service *service) GetEventsPending() ([]Event, error) {
	var events []Event

	err := service.db.Preload("Channel").Preload("Input").Preload("Channel.Outputs").Find(&events, "events.active = ? AND events.time <= ?", true, time.Now()).Error
	return events, err
}

func (service *service) UpdateEvent(event *Event) (*Event, error) {
	err := service.db.Save(event).Error

	return event, err
}

func (service *service) DeleteEvent(event *Event) error {
	return service.db.Delete(event).Error
}
