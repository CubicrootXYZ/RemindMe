package asyncmessenger

import "maunium.net/go/mautrix/event"

// MessageTypes available
var (
	messageTypeText     = "m.text"
	messageTypeReaction = "m.reaction"
)

// EventTypes available
var (
	eventTypeRoomMessage = "m.room.message"
)

// Formats available
var (
	formatCustomHTML = "org.matrix.custom.html"
)

// Relations available
var (
	relationAnnotiation = "m.annotation"
)

// Mimetypes available for MSC1767 events https://github.com/matrix-org/matrix-spec-proposals/blob/matthew/msc1767/proposals/1767-extensible-events.md
var (
	mimetypeTextPlain = "text/plain"
	mimetypeTextHTML  = "text/html"
)

type messageEvent struct {
	Body          string `json:"body,omitempty"`
	Format        string `json:"format,omitempty"`
	FormattedBody string `json:"formatted_body,omitempty"`
	MsgType       string `json:"msgtype,omitempty"`
	Type          string `json:"type,omitempty"`
	RelatesTo     struct {
		EventID   string `json:"event_id,omitempty"`
		Key       string `json:"key,omitempty"`
		RelType   string `json:"rel_type,omitempty"`
		InReplyTo *struct {
			EventID string `json:"event_id,omitempty"`
		} `json:"m.in_reply_to,omitempty"`
	} `json:"m.relates_to,omitempty"`
	MSC1767Message []matrixMSC1767Event `json:"org.matrix.msc1767.message,omitempty"`
}

func (messageEvent *messageEvent) getEventType() event.Type {
	if messageEvent.Type == messageTypeReaction {
		return event.EventReaction
	}

	return event.EventMessage
}

// MatrixMSC1767Message defines a MSC1767 message
type matrixMSC1767Event struct {
	Body     string `json:"body"`
	Mimetype string `json:"mimetype"`
}
