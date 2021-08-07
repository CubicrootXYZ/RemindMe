package database

import (
	"time"

	"gorm.io/gorm"
)

// Reminder is the database object for a reminder
type Reminder struct {
	gorm.Model
	RemindTime     time.Time
	Message        string
	Active         bool
	RepeatInterval uint64
	RepeatMax      uint64
	ChannelID      uint
	Channel        Channel
}

// GetPendingReminders returns a list with all pending reminders for the given channel
func (d *Database) GetPendingReminders(channel *Channel) ([]Reminder, error) {
	reminders := make([]Reminder, 0)

	err := d.db.Find(&reminders, "channel_id = ? AND active = ?", channel.ID, true).Error

	return reminders, err
}

// GetMessageFromReminder returns the message with the specified message type regarding the reminder
func (d *Database) GetMessageFromReminder(reminderID uint, msgType MessageType) (*Message, error) {
	message := &Message{}
	err := d.db.Find(message, "reminder_id = ? AND type = ?", reminderID, msgType).Error
	return message, err
}

// GetPendingReminder returns all reminders that are due
func (d *Database) GetPendingReminder() (*[]Reminder, error) {
	reminder := make([]Reminder, 0)
	err := d.db.Debug().Preload("Channel").Find(&reminder, "active = ? AND remind_time <= ?", 1, time.Now()).Error

	return &reminder, err
}

// SetReminderDone sets a reminder as inactive
func (d *Database) SetReminderDone(reminder *Reminder) (*Reminder, error) {
	reminder.Active = false
	err := d.db.Save(reminder).Error

	return reminder, err
}

// AddReminder adds a reminder to the database
func (d *Database) AddReminder(remindTime time.Time, message string, active bool, repeatInterval uint64, channel *Channel) (*Reminder, error) {
	reminder := Reminder{
		Message:        message,
		RemindTime:     remindTime,
		Active:         true,
		RepeatInterval: repeatInterval,
		RepeatMax:      0,
		Channel:        *channel,
		ChannelID:      channel.ID,
	}
	reminder.Model.CreatedAt = time.Now().UTC()

	err := d.db.Create(&reminder).Error

	return &reminder, err
}

// DeleteReminder deletes the given reminder
func (d *Database) DeleteReminder(reminderID uint) (*Reminder, error) {
	deleteReminder := &Reminder{}
	err := d.db.First(&deleteReminder, "id = ?", reminderID).Error
	if err != nil {
		return nil, err
	}

	err = d.db.Delete(&deleteReminder).Statement.Error
	return deleteReminder, err
}

// UpdateReminder updates the reminder
func (d *Database) UpdateReminder(reminderID uint, remindTime time.Time) (*Reminder, error) {
	reminder := &Reminder{}
	err := d.db.First(reminder, "id = ?", reminderID).Error
	if err != nil {
		return nil, err
	}

	reminder.RemindTime = remindTime
	reminder.Active = true
	err = d.db.Save(reminder).Error
	return reminder, err
}

// GetMessagesByReminderID returns a list with all messages for the given reminder id
func (d *Database) GetMessagesByReminderID(id uint) ([]*Message, error) {
	messages := make([]*Message, 0)
	err := d.db.Find(&messages, "reminder_id = ?", id).Error

	return messages, err
}
