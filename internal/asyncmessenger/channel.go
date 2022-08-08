package asyncmessenger

import (
	"time"

	"gorm.io/gorm"
)

// Channel holds all information about a channel
type Channel struct {
	gorm.Model
	Created           time.Time
	ChannelIdentifier string
	UserIdentifier    string
	TimeZone          string
	DailyReminder     *uint
	CalendarSecret    string
	LastCryptoEvent   string
}
