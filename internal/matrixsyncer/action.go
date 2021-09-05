package matrixsyncer

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix/event"
)

// Action defines an action the user can perform
type Action struct {
	Name     string   // Name of the action just for displaying
	Examples []string // Example commands to trigger the action
	Regex    string   // Regex the message must match to trigger the action
	Action   func(evt *event.Event, channel *database.Channel) error
}

// ReplyAction defines actions that are performed on replies
type ReplyAction struct {
	Name         string                 // Name of the action just for displaying
	Examples     []string               // Example commands to trigger the action
	Regex        string                 // Regex the message must match to trigger the action
	ReplyToTypes []database.MessageType // Kind of message the reply is for
	Action       func(evt *event.Event, channel *database.Channel, replyMessage *database.Message, content *event.MessageEventContent) error
}

// ReactionAction defines an action performed on receiving a reaction
type ReactionAction struct {
	Name   string   // Name of the action just for displaying
	Keys   []string // The key the reaction must match
	Type   ReactionActionType
	Action func(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error
}

// ReactionActionType defines types of reaction actions
type ReactionActionType string

const (
	ReactionActionTypeReminderRequest = ReactionActionType(string(database.MessageTypeReminderRequest))
	ReactionActionTypeReminder        = ReactionActionType(string(database.MessageTypeReminder))
	ReactionActionTypeReminderSuccess = ReactionActionType(string(database.MessageTypeReminderSuccess))
	ReactionActionTypeDailyReminder   = ReactionActionType(string(database.MessageTypeDailyReminder))
	ReactionActionTypeAll             = ReactionActionType("")
)
