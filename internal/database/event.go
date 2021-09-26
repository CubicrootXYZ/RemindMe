package database

import "gorm.io/gorm"

// Event holds information about a matrix event
type Event struct {
	gorm.Model
	ChannelID          uint `gorm:"index"`
	Channel            Channel
	Timestamp          int64
	ExternalIdentifier string    `gorm:"uniqueIndex;size:500"`
	EventType          EventType `gorm:"index;size:250"`
	EventSubType       string    `gorm:"index;size:250"`
	AdditionalInfo     string
}

// EventType is the type of the event
type EventType string

const (
	EventTypeMembership = EventType("MEMBERSHIP")
)

// GET DATA

// IsEventKnown returns if the given externalID of an event is already registered
func (d *Database) IsEventKnown(externalID string) (bool, error) {
	var evt Event
	err := d.db.First(&evt, "external_identifier = ?", externalID).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// INSERT DATA

// AddEvent adds the given event to the database
func (d *Database) AddEvent(event *Event) (*Event, error) {
	err := d.db.Create(event).Error
	return event, err
}
