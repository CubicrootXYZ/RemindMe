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
