package database

import (
	"time"

	"gorm.io/gorm"
	"maunium.net/go/mautrix/event"
)

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
	MessageTypeReminderRequest = MessageType("REMINDER_REQUEST")
	MessageTypeReminderSuccess = MessageType("REMINDER_SUCCESS")
	MessageTypeReminder        = MessageType("REMINDER")
	MessageTypeReminderUpdate  = MessageType("REMINDER_UPDATE")
	MessageTypeReminderDelete  = MessageType("REMINDER_DELETE")
)

// AddMessageFromMatrix adds a message to the database
func (d *Database) AddMessageFromMatrix(id string, timestamp int64, content *event.MessageEventContent, reminder *Reminder, msgType MessageType, channel *Channel) (*Message, error) {
	relatesTo := ""
	if content.RelatesTo != nil {
		relatesTo = content.RelatesTo.EventID.String()
	}
	message := Message{
		Body:               content.Body,
		BodyHTML:           content.FormattedBody,
		ReminderID:         reminder.ID,
		Reminder:           *reminder,
		ResponseToMessage:  relatesTo,
		Type:               msgType,
		ChannelID:          channel.ID,
		Channel:            *channel,
		Timestamp:          timestamp,
		ExternalIdentifier: id,
	}
	message.Model.CreatedAt = time.Now().UTC()

	err := d.db.Create(&message).Error

	return &message, err
}

// AddMessage adds a message to the database
func (d *Database) AddMessage(message *Message) (*Message, error) {

	err := d.db.Create(message).Error

	return message, err
}

// GetMessageByExternalID returns if found the message with the given external id
func (d *Database) GetMessageByExternalID(externalID string) (*Message, error) {
	message := &Message{}
	err := d.db.Preload("Reminder").First(&message, "external_identifier = ?", externalID).Error
	return message, err
}
