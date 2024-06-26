package tests

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
)

type MatrixMessageAssertion func(*matrixdb.MatrixMessage)

func ExpectNewMessageFromEvent(matrixDB *matrixdb.MockService, event *matrix.MessageEvent, t matrixdb.MatrixMessageType, assertions ...MatrixMessageAssertion) {
	sender := event.Event.Sender.String()
	msg := &matrixdb.MatrixMessage{
		ID:            event.Event.ID.String(),
		UserID:        &sender,
		RoomID:        event.Room.ID,
		Body:          event.Content.Body,
		BodyFormatted: event.Content.FormattedBody,
		SendAt:        time.UnixMilli(event.Event.Timestamp),
		Incoming:      true,
		Type:          t,
	}

	for _, assertion := range assertions {
		assertion(msg)
	}

	matrixDB.EXPECT().NewMessage(msg).Return(nil, nil)
}

func MsgWithDBEventID(id uint) MatrixMessageAssertion {
	return func(mm *matrixdb.MatrixMessage) {
		mm.EventID = &id
	}
}
