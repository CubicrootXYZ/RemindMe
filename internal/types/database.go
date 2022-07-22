package types

import (
	"database/sql"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"maunium.net/go/mautrix/event"
)

// Database defines an interface for a data storage provider
type Database interface {
	// Reminders
	GetPendingReminders(channel *database.Channel) ([]database.Reminder, error)
	GetReminderForChannelIDByID(channelID string, reminderID int) (*database.Reminder, error)

	AddReminder(remindTime time.Time, message string, active bool, repeatInterval uint64, channel *database.Channel) (*database.Reminder, error)

	UpdateReminder(reminderID uint, remindTime time.Time, repeatInterval uint64, repeatTimes uint64) (*database.Reminder, error)

	DeleteReminder(reminderID uint) (*database.Reminder, error)
	// Messages
	AddMessage(message *database.Message) (*database.Message, error)

	GetMessageByExternalID(externalID string) (*database.Message, error)
	GetMessagesByReminderID(id uint) ([]*database.Message, error)
	GetLastMessageByTypeForReminder(msgType database.MessageType, reminderID uint) (*database.Message, error)

	AddMessageFromMatrix(id string, timestamp int64, content *event.MessageEventContent, reminder *database.Reminder, msgType database.MessageType, channel *database.Channel) (*database.Message, error)

	// Channels
	GetChannel(id uint) (*database.Channel, error)
	GetChannelByUserIdentifier(userID string) (*database.Channel, error)
	GetChannelsByUserIdentifier(userID string) ([]database.Channel, error)
	GetChannelsByChannelIdentifier(channelID string) ([]database.Channel, error)
	GetChannelByUserAndChannelIdentifier(userID string, channelID string) (*database.Channel, error)
	GetChannelList() ([]database.Channel, error)
	ChannelCount() (int64, error)

	GenerateNewCalendarSecret(channel *database.Channel) error
	UpdateChannel(channelID uint, timeZone string, dailyReminder *uint, role *roles.Role) (*database.Channel, error)
	ChannelSaveChanges(channel *database.Channel) error

	AddChannel(userID, channelID string, role roles.Role) (*database.Channel, error)

	CleanAdminChannels(keep []*database.Channel) error
	DeleteChannel(channel *database.Channel) error
	DeleteChannelsFromUser(userID string) error

	// Events
	IsEventKnown(externalID string) (bool, error)

	AddEvent(event *database.Event) (*database.Event, error)

	// Blocklist
	IsUserBlocked(userID string) (bool, error)
	GetBlockedUserList() ([]database.Blocklist, error)

	AddUserToBlocklist(userID string, reason string) error

	RemoveUserFromBlocklist(userID string) error

	SQLDB() (*sql.DB, error)
}
