package types

import (
	"maunium.net/go/mautrix/event"
)

// MessageEvent holds data from a message event
type MessageEvent struct {
	Event       *event.Event
	Content     *event.MessageEventContent
	IsEncrypted bool
}
