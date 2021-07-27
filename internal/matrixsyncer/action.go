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
