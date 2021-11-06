package matrixsyncer

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix"
)

// MatrixMessage holds information for a matrix response message
type MatrixMessage struct {
	Body          string `json:"body"`
	Format        string `json:"format"`
	FormattedBody string `json:"formatted_body,omitempty"`
	MsgType       string `json:"msgtype"`
	Type          string `json:"type"`
	Relatesto     struct {
		InReplyTo struct {
			EventID string `json:"event_id,omitempty"`
		} `json:"m.in_reply_to,omitempty"`
	} `json:"m.relates_to,omitempty"`
}

// Messenger defines an interface for interacting with matrix messages
type Messenger interface {
	SendReplyToEvent(msg string, replyEvent *types.MessageEvent, channel *database.Channel, msgType database.MessageType) (resp *mautrix.RespSendEvent, err error)
	CreateChannel(userID string) (*mautrix.RespCreateRoom, error)
	SendFormattedMessage(msg, msgFormatted string, channel *database.Channel, msgType database.MessageType, relatedReminderID uint) (resp *mautrix.RespSendEvent, err error)
	DeleteMessage(messageID, roomID string) error
	SendNotice(msg, roomID string) (resp *mautrix.RespSendEvent, err error)
}
