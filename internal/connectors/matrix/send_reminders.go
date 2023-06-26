package matrix

import (
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
)

func (service *service) SendReminder(event *daemon.Event, output *daemon.Output) error {
	room, err := service.matrixDatabase.GetRoomByID(output.OutputID)
	if err != nil {
		return err
	}

	originalMessage, err := service.matrixDatabase.GetEventMessageByOutputAndEvent(event.ID, output.OutputID, output.OutputType)
	if err != nil {
		service.logger.Err(err)
		originalMessage = nil
	}

	message, messageFormatted, err := format.MessageFromEvent(event, "") // TODO get timezone
	if err != nil {
		return err
	}

	var resp *messenger.MessageResponse
	if originalMessage == nil {
		resp, err = service.messenger.SendMessage(messenger.HTMLMessage(
			message,
			messageFormatted,
			room.RoomID,
		))
		if err != nil {
			return err
		}
	} else {
		resp, err = service.messenger.SendResponse(&messenger.Response{
			Message:                   message,
			MessageFormatted:          messageFormatted,
			RespondToMessage:          originalMessage.Body,
			RespondToMessageFormatted: originalMessage.BodyFormatted,
			RespondToEventID:          originalMessage.ID,
			ChannelExternalIdentifier: room.RoomID,
		})
		if err != nil {
			return err
		}
	}

	dbMsg := &matrixdb.MatrixMessage{
		ID:            resp.ExternalIdentifier,
		UserID:        nil, // There is no user, events can be from any source
		RoomID:        room.ID,
		Body:          message,
		BodyFormatted: messageFormatted,
		SendAt:        resp.Timestamp,
		Type:          matrixdb.MessageTypeEvent,
		EventID:       &event.ID,
	}
	if originalMessage != nil {
		dbMsg.ReplyToMessageID = &originalMessage.ID
	}
	_, err = service.matrixDatabase.NewMessage(dbMsg)
	if err != nil {
		service.logger.Errorf("failed to save message to database: %v", err)
	}

	return nil
}

func (service *service) SendDailyReminder(*daemon.DailyReminder, *daemon.Output) error {
	return nil
}