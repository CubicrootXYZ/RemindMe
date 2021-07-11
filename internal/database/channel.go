package database

import (
	"time"

	"gorm.io/gorm"
)

// Channel holds data about a messaging channel
type Channel struct {
	gorm.Model
	Created           time.Time
	ChannelIdentifier string
	UserIdentifier    string
}
