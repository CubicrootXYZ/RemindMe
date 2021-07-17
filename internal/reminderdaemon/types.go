package reminderdaemon

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

type Database interface {
	GetPendingReminder() (*[]database.Reminder, error)
	AddMessage(message *database.Message) (*database.Message, error)
	GetMessageFromReminder(reminderID uint, msgType database.MessageType) (*database.Message, error)
	SetReminderDone(*database.Reminder) (*database.Reminder, error)
}

type Messenger interface {
	SendReminder(*database.Reminder, *database.Message) (*database.Message, error)
}
