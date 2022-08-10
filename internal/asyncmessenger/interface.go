package asyncmessenger

// Messenger provides an interface for sending and handling matrix events
// The Queue... methods work asynchronous but do not provide any information about the sent message
// The Send... methods work snchronous and provide more detailed feedback, they can be used asynchronous
type Messenger interface {
	QueueMessage(message *Message) error
	SendMessage(message *Message) (*MessageResponse, error)
	QueueReaction(reaction *Reaction) error
	QueueResponse(response *Response) error
	SendResponse(response *Response) (*MessageResponse, error)
	QueueRedact(redact *Redact) error
	QueueEdit(edit *Edit) error
}
