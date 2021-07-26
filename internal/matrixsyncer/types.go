package matrixsyncer

import (
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
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

type Messenger interface {
	SendReplyToEvent(msg string, replyEvent *event.Event, roomID string) (resp *mautrix.RespSendEvent, err error)
	CreateChannel(userID string) (*mautrix.RespCreateRoom, error)
	SendFormattedMessage(msg, msgFormatted, roomID string) (resp *mautrix.RespSendEvent, err error)
}
