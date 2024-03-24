package tests

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type ReactionEventOpt func(evt *matrix.ReactionEvent)

func TestReactionEvent(opts ...ReactionEventOpt) *matrix.ReactionEvent {
	evt := &matrix.ReactionEvent{
		Event: &event.Event{
			ID:        id.EventID("evt1"),
			Sender:    id.UserID("@user:example.com"),
			RoomID:    id.RoomID("!room123"),
			Timestamp: 92848488,
		},
		Content: &event.ReactionEventContent{
			RelatesTo: event.RelatesTo{
				Key:     "X",
				EventID: id.EventID("msg1"),
			},
		},
		Room: &matrixdb.MatrixRoom{
			RoomID:   "!room123",
			Users:    []matrixdb.MatrixUser{},
			TimeZone: "UTC",
		},
		Channel: &database.Channel{
			Model: gorm.Model{
				ID: 68272,
			},
		},
		Input: &database.Input{
			Model: gorm.Model{
				ID: 187,
			},
		},
	}

	for _, opt := range opts {
		opt(evt)
	}

	return evt
}

func ReactionWithKey(key string) ReactionEventOpt {
	return func(evt *matrix.ReactionEvent) {
		evt.Content.RelatesTo.Key = key
	}
}
