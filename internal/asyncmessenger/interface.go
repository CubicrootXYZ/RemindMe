package asyncmessenger

import "errors"

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
	SendEditAsync(edit *Edit) error
}

// Errors returned by the messenger
var (
	ErrRetriesExceeded = errors.New("amount of retries exceeded")
)
