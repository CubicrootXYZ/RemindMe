package mapping

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
)

// MessageFromEvent creates a MatrixMessage from a MessageEvent.
func MessageFromEvent(event *matrix.MessageEvent) *database.MatrixMessage {
	sender := event.Event.Sender.String()

	return &database.MatrixMessage{
		ID:            event.Event.ID.String(),
		UserID:        &sender,
		RoomID:        event.Room.ID,
		Body:          event.Content.Body,
		BodyFormatted: event.Content.FormattedBody,
		SendAt:        time.UnixMilli(event.Event.Timestamp),
		Incoming:      true,
	}
}
