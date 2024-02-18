package tests

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type EventOpt func(evt *matrix.MessageEvent)

func TestEvent(opts ...EventOpt) *matrix.MessageEvent {
	evt := &matrix.MessageEvent{
		Event: &event.Event{
			ID:        id.EventID("evt1"),
			Sender:    id.UserID("@user:example.com"),
			RoomID:    id.RoomID("!room123"),
			Timestamp: 92848488,
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

func WithBody(body, formattedBody string) EventOpt {
	return func(evt *matrix.MessageEvent) {
		evt.Content.Body = body
		evt.Content.FormattedBody = formattedBody
	}
}

func WithUserInRoom(user matrixdb.MatrixUser) EventOpt {
	return func(evt *matrix.MessageEvent) {
		evt.Room.Users = append(evt.Room.Users, user)
	}
}

func WithOutput(output database.Output) EventOpt {
	return func(evt *matrix.MessageEvent) {
		if evt.Channel.Outputs == nil {
			evt.Channel.Outputs = make([]database.Output, 0)
		}
		evt.Channel.Outputs = append(evt.Channel.Outputs, output)
	}
}

type MessageOpt func(msg *matrixdb.MatrixMessage)

func TestMessage(opts ...MessageOpt) *matrixdb.MatrixMessage {
	msg := &matrixdb.MatrixMessage{
		ID:      "msg1",
		EventID: ToP(uint(1)),
		Event: &database.Event{
			Model: gorm.Model{
				ID: 1,
			},
			Message: "test event",
		},
	}

	for _, o := range opts {
		o(msg)
	}
	return msg
}

func WithFromTestEvent() MessageOpt {
	return func(msg *matrixdb.MatrixMessage) {
		msg.ID = "evt1"
		msg.UserID = ToP("@user:example.com")
		msg.Body = "message"
		msg.BodyFormatted = "<b>message</b>"
		msg.SendAt = time.UnixMilli(92848488)
		msg.Incoming = true
		msg.Event = nil
		msg.EventID = nil
	}
}

func WithMessageType(mt matrixdb.MatrixMessageType) MessageOpt {
	return func(msg *matrixdb.MatrixMessage) {
		msg.Type = mt
	}
}
