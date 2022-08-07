package types

import (
	"regexp"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix/event"
)

// ReplyAction defines actions that are performed on replies
type ReplyAction struct {
	Name         string                 // Name of the action just for displaying
	Examples     []string               // Example commands to trigger the action
	Regex        *regexp.Regexp         // Regex the message must match to trigger the action
	ReplyToTypes []database.MessageType // Kind of message the reply is for
	Action       func(evt *MessageEvent, channel *database.Channel, replyMessage *database.Message) error
}

// Action defines an action the user can perform
type Action struct {
	Name     string         // Name of the action just for displaying
	Examples []string       // Example commands to trigger the action
	Regex    *regexp.Regexp // Regex the message must match to trigger the action
	Action   func(evt *MessageEvent, channel *database.Channel) error
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

func (reactionActionType *ReactionActionType) HumanReadable() string {
	switch *reactionActionType {
	case ReactionActionTypeReminderRequest:
		return "reminder request"
	case ReactionActionTypeReminder:
		return "reminder message"
	case ReactionActionTypeReminderSuccess:
		return "reminder confirmation"
	case ReactionActionTypeDailyReminder:
		return "daily reminder"
	case ReactionActionTypeAll:
		return "all messages"
	default:
		return "unknown message"
	}
}

// ReactionAction defines an action performed on receiving a reaction
type ReactionAction struct {
	Name   string   // Name of the action just for displaying
	Keys   []string // The key the reaction must match
	Type   ReactionActionType
	Action func(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error
}
