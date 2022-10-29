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
	ReminderID         *uint `gorm:"index"`
	Reminder           Reminder
	ResponseToMessage  string `gorm:"index"` // External identifier of parent message
	Type               MessageType
	ChannelID          uint `gorm:"index"`
	Channel            Channel
	Timestamp          int64
	ExternalIdentifier string `gorm:"index"`
}

// MessageType defines different types of messages
type MessageType string

// Message types differentiate the context of a message
const (
	// Reminder itself
	MessageTypeReminderRequest = MessageType("REMINDER_REQUEST")
	MessageTypeReminderSuccess = MessageType("REMINDER_SUCCESS")
	MessageTypeReminderFail    = MessageType("REMINDER_FAIL")
	MessageTypeReminder        = MessageType("REMINDER")
	// Arbitrary actions
	MessageTypeActions          = MessageType("ACTIONS")
	MessageTypeReminderList     = MessageType("REMINDER_LIST")
	MessageTypeIcalLink         = MessageType("ICAL_LINK")
	MessageTypeIcalLinkRequest  = MessageType("ICAL_LINK_REQUEST")
	MessageTypeIcalRenew        = MessageType("ICAL_RENEW")
	MessageTypeIcalRenewRequest = MessageType("ICAL_RENEW_REQUEST")
	MessageTypeWelcome          = MessageType("WELCOME")
	// Reminder edits
	MessageTypeReminderUpdate           = MessageType("REMINDER_UPDATE")
	MessageTypeReminderUpdateFail       = MessageType("REMINDER_UPDATE_FAIL")
	MessageTypeReminderUpdateSuccess    = MessageType("REMINDER_UPDATE_SUCCESS")
	MessageTypeReminderDelete           = MessageType("REMINDER_DELETE")
	MessageTypeReminderDeleteSuccess    = MessageType("REMINDER_DELETE_SUCCESS")
	MessageTypeReminderDeleteFail       = MessageType("REMINDER_DELETE_Fail")
	MessageTypeReminderRecurringRequest = MessageType("REMINDER_RECURRING_REQUEST")
	MessageTypeReminderRecurringSuccess = MessageType("REMINDER_RECURRING_SUCCESS")
	MessageTypeReminderRecurringFail    = MessageType("REMINDER_RECURRING_FAIL")
	// Settings
	MessageTypeTimezoneChangeRequest        = MessageType("TIMEZONE_CHANGE")
	MessageTypeTimezoneChangeRequestSuccess = MessageType("TIMEZONE_CHANGE_SUCCESS")
	MessageTypeTimezoneChangeRequestFail    = MessageType("TIMEZONE_CHANGE_FAIL")
	// Daily Reminder
	MessageTypeDailyReminder              = MessageType("DAILY_REMINDER")
	MessageTypeDailyReminderUpdate        = MessageType("DAILY_REMINDER_UPDATE")
	MessageTypeDailyReminderUpdateFail    = MessageType("DAILY_REMINDER_UPDATE_FAIL")
	MessageTypeDailyReminderUpdateSuccess = MessageType("DAILY_REMINDER_UPDATE_SUCCESS")
	MessageTypeDailyReminderDelete        = MessageType("DAILY_REMINDER_DELETE")
	MessageTypeDailyReminderDeleteFail    = MessageType("DAILY_REMINDER_DELETE_FAIL")
	MessageTypeDailyReminderDeleteSuccess = MessageType("DAILY_REMINDER_DELETE_SUCCESS")
	// Do not save!
	MessageTypeDoNotSave = MessageType("")
)

// MessageTypesWithReminder message types with reminders
var MessageTypesWithReminder = []MessageType{MessageTypeReminderRequest, MessageTypeReminderSuccess, MessageTypeReminderUpdate, MessageTypeReminderUpdateSuccess, MessageTypeReminderRecurringRequest, MessageTypeReminderRecurringSuccess, MessageTypeReminderRecurringFail, MessageTypeReminderUpdateFail}

// GET DATA

// GetMessageByExternalID returns if found the message with the given external id
func (d *Database) GetMessageByExternalID(externalID string) (*Message, error) {
	message := &Message{}
	err := d.db.Preload("Reminder").Preload("Channel").First(&message, "external_identifier = ?", externalID).Error
	return message, err
}

// GetMessagesByReminderID returns a list with all messages for the given reminder id
func (d *Database) GetMessagesByReminderID(id uint) ([]*Message, error) {
	messages := make([]*Message, 0)
	err := d.db.Find(&messages, "reminder_id = ?", id).Error

	return messages, err
}

// GetLastMessageByType returns the last message in the given channel with the given message type
func (d *Database) GetLastMessageByType(msgType MessageType, channel *Channel) (*Message, error) {
	message := &Message{}
	err := d.db.Order("timestamp desc").First(message, "channel_id = ? AND type = ?", channel.ID, msgType).Error

	return message, err
}

// GetLastMessageByTypeForReminder returns the last message of the specified type tied to the given reminder id
func (d *Database) GetLastMessageByTypeForReminder(msgType MessageType, reminderID uint) (*Message, error) {
	message := &Message{}
	err := d.db.Joins("Channel").Order("timestamp desc").First(message, "reminder_id = ? AND type = ?", reminderID, msgType).Error

	return message, err
}

// INSERT DATA

// AddMessageFromMatrix adds a message to the database
func (d *Database) AddMessageFromMatrix(id string, timestamp int64, content *event.MessageEventContent, reminder *Reminder, msgType MessageType, channel *Channel) (*Message, error) {
	relatesTo := ""
	if content.RelatesTo != nil {
		relatesTo = content.RelatesTo.EventID.String()
	}
	message := Message{
		ResponseToMessage:  relatesTo,
		Type:               msgType,
		Timestamp:          timestamp,
		ExternalIdentifier: id,
	}

	if channel != nil {
		message.ChannelID = channel.ID
		message.Channel = *channel
	}
	message.Model.CreatedAt = time.Now().UTC()

	if content != nil {
		message.Body = content.Body
		message.BodyHTML = content.FormattedBody
	}

	if reminder != nil {
		message.Reminder = *reminder
		message.ReminderID = &reminder.ID
	}

	err := d.db.Create(&message).Error

	return &message, err
}

// AddMessage adds a message to the database
func (d *Database) AddMessage(message *Message) (*Message, error) {
	err := d.db.Create(message).Error

	return message, err
}
