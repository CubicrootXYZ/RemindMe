package asyncmessenger

// Message holds information about a message
type Message struct {
	Body              string
	BodyHTML          string
	ReminderID        *uint
	Reminder          Reminder
	ResponseToMessage string
	ChannelID         uint
	Channel           Channel
}

type MessageResponse struct {
	ExternalIdentifier string
	Timestamp          int64
}
