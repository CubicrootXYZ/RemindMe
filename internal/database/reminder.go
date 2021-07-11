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
