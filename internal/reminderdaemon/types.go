package reminderdaemon

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix"
)

// Database defines a database interface for the reminderdaemon
type Database interface {
	// Message
	AddMessage(message *database.Message) (*database.Message, error)
	GetLastMessageByType(msgType database.MessageType, channel *database.Channel) (*database.Message, error)
	// Reminder
	GetPendingReminder() ([]database.Reminder, error)
	GetMessageFromReminder(reminderID uint, msgType database.MessageType) (*database.Message, error)
	SetReminderDone(*database.Reminder) (*database.Reminder, error)
	GetDailyReminder(channel *database.Channel) (*[]database.Reminder, error)
	// Channel
	GetChannelList() ([]database.Channel, error)
}

// Messenger defines a messenger interface for the reminderdaemon
type Messenger interface {
	SendReminder(*database.Reminder, *database.Message) (*database.Message, error)
	SendFormattedMessage(msg, msgFormatted string, channel *database.Channel, msgType database.MessageType, relatedReminderID uint) (resp *mautrix.RespSendEvent, err error)
}
