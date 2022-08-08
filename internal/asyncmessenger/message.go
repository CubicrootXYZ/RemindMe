package asyncmessenger

import "gorm.io/gorm"

// Message holds information about a message
type Message struct {
	gorm.Model
	Body               string
	BodyHTML           string
	ReminderID         *uint
	Reminder           Reminder
	ResponseToMessage  string
	Type               string
	ChannelID          uint
	Channel            Channel
	Timestamp          int64
	ExternalIdentifier string
}
