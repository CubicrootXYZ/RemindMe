package types

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix"
)

// Messenger defines an interface for interacting with matrix messages
type Messenger interface {
	SendReplyToEvent(msg string, replyEvent *MessageEvent, channel *database.Channel, msgType database.MessageType) (resp *mautrix.RespSendEvent, err error)
	CreateChannel(userID string) (*mautrix.RespCreateRoom, error)
	SendFormattedMessage(msg, msgFormatted string, channel *database.Channel, msgType database.MessageType, relatedReminderID uint) (resp *mautrix.RespSendEvent, err error)
	DeleteMessage(messageID, roomID string) error
	SendNotice(msg, roomID string) (resp *mautrix.RespSendEvent, err error)
	SendReaction(reaction string, toMessage string, channel *database.Channel) (resp *mautrix.RespSendEvent, err error)
}
