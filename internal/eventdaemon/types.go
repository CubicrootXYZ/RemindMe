package eventdaemon

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix/event"
)

// Database defines an interface for a database
type Database interface {
	AddReminder(remindTime time.Time, message string, active bool, repeatInterval uint64, channel *database.Channel) (*database.Reminder, error)
	AddMessageFromMatrix(id string, timestamp int64, content *event.MessageEventContent, reminder *database.Reminder, msgType database.MessageType, channel *database.Channel) (*database.Message, error)
	GetChannelByUserIdentifier(userID string) (*database.Channel, error)
	GetChannelByUserAndChannelIdentifier(userID string, channelID string) (*database.Channel, error)
	AddChannel(userID, channelID string) (*database.Channel, error)
}

// Syncer is responsible for receiving messages from a messenger
type Syncer interface {
	Start(daemon *Daemon) error
	Stop()
}