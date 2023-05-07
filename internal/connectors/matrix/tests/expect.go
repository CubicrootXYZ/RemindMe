package tests

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
)

func ExpectNewMessageFromEvent(matrixDB *matrixdb.MockService, event *matrix.MessageEvent, t matrixdb.MatrixMessageType) {
	sender := event.Event.Sender.String()
	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            event.Event.ID.String(),
		UserID:        &sender,
		RoomID:        event.Room.ID,
		Body:          event.Content.Body,
		BodyFormatted: event.Content.FormattedBody,
		SendAt:        time.UnixMilli(event.Event.Timestamp),
		Incoming:      true,
		Type:          t,
	}).Return(nil, nil)
}
