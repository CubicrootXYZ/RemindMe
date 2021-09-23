package types

import "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"

// Database defines an interface for a data storage provider
type Database interface {
	// Channels
	GetChannel(id uint) (*database.Channel, error)
	GetChannelList() ([]database.Channel, error)
	GenerateNewCalendarSecret(channel *database.Channel) error
	// Reminders
	GetPendingReminders(channel *database.Channel) ([]database.Reminder, error)
}
