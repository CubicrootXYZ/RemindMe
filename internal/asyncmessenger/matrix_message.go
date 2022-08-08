package asyncmessenger

// MessageTypes available
var (
	MessageTypeText = "m.text"
)

// EventTypes available
var (
	EventTypeRoomMessage = "m.room.message"
)

// Formats available
var (
	FormatCustomHTML = "org.matrix.custom.html"
)

// MatrixMessage holds information for a matrix response message
type MatrixMessage struct {
	Body          string `json:"body,omitempty"`
	Format        string `json:"format,omitempty"`
	FormattedBody string `json:"formatted_body,omitempty"`
	MsgType       string `json:"msgtype,omitempty"`
	Type          string `json:"type,omitempty"`
	RelatesTo     struct {
		EventID   string `json:"event_id,omitempty"`
		Key       string `json:"key,omitempty"`
		RelType   string `json:"rel_type,omitempty"`
		InReplyTo struct {
			EventID string `json:"event_id,omitempty"`
		} `json:"m.in_reply_to,omitempty"`
	} `json:"m.relates_to,omitempty"`
	MSC1767Message []MatrixMSC1767Message `json:"org.matrix.msc1767.message,omitempty"`
}

// MatrixMSC1767Message defines a MSC1767 message
type MatrixMSC1767Message struct {
	Body     string `json:"body"`
	Mimetype string `json:"mimetype"`
}
