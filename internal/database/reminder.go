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
	Repeated       *uint64
	ChannelID      uint
	Channel        Channel
}

// GET DATA

// GetPendingReminders returns a list with all pending reminders for the given channel
func (d *Database) GetPendingReminders(channel *Channel) ([]Reminder, error) {
	reminders := make([]Reminder, 0)

	err := d.db.Find(&reminders, "channel_id = ? AND active = ?", channel.ID, true).Error

	return reminders, err
}

// GetMessageFromReminder returns the message with the specified message type regarding the reminder
func (d *Database) GetMessageFromReminder(reminderID uint, msgType MessageType) (*Message, error) {
	message := &Message{}
	err := d.db.Preload("Channel").Find(message, "reminder_id = ? AND type = ?", reminderID, msgType).Error
	return message, err
}

// GetPendingReminder returns all reminders that are due
func (d *Database) GetPendingReminder() (*[]Reminder, error) {
	reminder := make([]Reminder, 0)
	err := d.db.Debug().Preload("Channel").Find(&reminder, "active = ? AND remind_time <= ?", 1, time.Now().UTC()).Error

	return &reminder, err
}

// GetDailyReminder returns the reminders alerting in the next 24 hours
func (d *Database) GetDailyReminder(channel *Channel) (*[]Reminder, error) {
	reminder := make([]Reminder, 0)
	err := d.db.Order("remind_time asc").Find(&reminder, "channel_id = ? AND remind_time <= ? AND active = ?", channel.ID, time.Now().Add(time.Hour*24), true).Error

	return &reminder, err
}

// UPDATE DATA

// SetReminderDone sets a reminder as inactive
func (d *Database) SetReminderDone(reminder *Reminder) (*Reminder, error) {
	if reminder.RepeatMax > *reminder.Repeated && reminder.RepeatInterval > 0 {
		reminder.RemindTime.Add(time.Duration(reminder.RepeatInterval) * time.Minute)
	} else {
		reminder.Active = false
	}

	*reminder.Repeated++

	err := d.db.Save(reminder).Error

	return reminder, err
}

// UpdateReminder updates the reminder
func (d *Database) UpdateReminder(reminderID uint, remindTime time.Time, repeatInterval uint64, repeatTimes uint64) (*Reminder, error) {
	reminder := &Reminder{}
	err := d.db.First(reminder, "id = ?", reminderID).Error
	if err != nil {
		return nil, err
	}

	if time.Until(remindTime) > 0 {
		reminder.RemindTime = remindTime
	}

	reminder.Active = true
	reminder.RepeatInterval = repeatInterval
	reminder.RepeatMax = repeatTimes
	err = d.db.Save(reminder).Error
	return reminder, err
}

// INSERT DATA

// AddReminder adds a reminder to the database
func (d *Database) AddReminder(remindTime time.Time, message string, active bool, repeatInterval uint64, channel *Channel) (*Reminder, error) {
	reminder := Reminder{
		Message:        message,
		RemindTime:     remindTime.UTC(),
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

// DELETE DATA

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
