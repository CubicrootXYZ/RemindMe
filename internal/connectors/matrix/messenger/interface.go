package messenger

import (
	"errors"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Messenger provides an interface for sending and handling matrix events
// The Queue... methods work asynchronous but do not provide any information about the sent message
// The Send... methods work snchronous and provide more detailed feedback, they can be used asynchronous
type Messenger interface {
	SendMessageAsync(message *Message) error
	SendMessage(message *Message) (*MessageResponse, error)
	SendReactionAsync(reaction *Reaction) error
	SendResponseAsync(response *Response) error
	SendResponse(response *Response) (*MessageResponse, error)
	SendRedactAsync(redact *Redact) error
	CreateChannel(userID string) (*ChannelResponse, error)
	DeleteMessageAsync(deleteAction *Delete) error
	// TODO future improvement: make room member cache flushable throug this interface and flush it on room member updates
}

// Errors returned by the messenger
var (
	ErrRetriesExceeded = errors.New("amount of retries exceeded")
)

// MatrixClient defines an interface to wrap the matrix API
type MatrixClient interface {
	SendMessageEvent(roomID id.RoomID, eventType event.Type, contentJSON interface{}, extra ...mautrix.ReqSendEvent) (resp *mautrix.RespSendEvent, err error)
	RedactEvent(roomID id.RoomID, eventID id.EventID, extra ...mautrix.ReqRedact) (resp *mautrix.RespSendEvent, err error)
	JoinedMembers(roomID id.RoomID) (resp *mautrix.RespJoinedMembers, err error)
	CreateRoom(req *mautrix.ReqCreateRoom) (resp *mautrix.RespCreateRoom, err error)
}
