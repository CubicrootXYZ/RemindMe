package message_test

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type eventOpt func(evt *matrix.MessageEvent)

func testEvent(opts ...eventOpt) *matrix.MessageEvent {
	evt := &matrix.MessageEvent{
		Event: &event.Event{
			ID:        id.EventID("evt1"),
			Sender:    id.UserID("@user:example.com"),
			RoomID:    id.RoomID("!room123"),
			Timestamp: 928484888888000000,
		},
		Content: &event.MessageEventContent{
			Body:          "message",
			FormattedBody: "<b>message</b>",
		},
		Room: &matrixdb.MatrixRoom{
			RoomID: "!room123",
			Users:  []matrixdb.MatrixUser{},
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

func withBody(body, formattedBody string) eventOpt {
	return func(evt *matrix.MessageEvent) {
		evt.Content.Body = body
		evt.Content.FormattedBody = formattedBody
	}
}

func withUserInRoom(user matrixdb.MatrixUser) eventOpt {
	return func(evt *matrix.MessageEvent) {
		evt.Room.Users = append(evt.Room.Users, user)
	}
}
