package database

import "gorm.io/gorm"

// Message holds information about a single message
type Message struct {
	gorm.Model
	Body               string
	BodyHTML           string
	ReminderID         uint
	Reminder           Reminder
	ResponseToMessage  string
	Type               MessageType
	ChannelID          uint
	Channel            Channel
	Timestamp          int64
	ExternalIdentifier string
}

// MessageType defines different types of messages
type MessageType string

const (
	MessageTypeReminderRequest = MessageType("REMINDER")
	MessageTypeReminderSuccess = MessageType("REMINDER_SUCCESS")
)
