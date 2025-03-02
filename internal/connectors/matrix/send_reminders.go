package matrix

import (
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
)

var ReminderReactions = []string{"✅", "▶️", "⏩", "1️⃣", "4️⃣"}
var ReminderReactionsRecurring = []string{"🔂"}

func (service *service) SendReminder(event *daemon.Event, output *daemon.Output) error {
	room, err := service.matrixDatabase.GetRoomByID(output.OutputID)
	if err != nil {
		return err
	}

	originalMessage, err := service.matrixDatabase.GetEventMessageByOutputAndEvent(event.ID, output.OutputID, output.OutputType)
	if err != nil {
		service.logger.Error("failed to get event message", "error", err)
		originalMessage = nil
	}

	message, messageFormatted, err := format.MessageFromEvent(event, room.TimeZone)
	if err != nil {
		return err
	}

	resp, err := service.messenger.SendMessage(messenger.HTMLMessage(
		message,
		messageFormatted,
		room.RoomID,
	))
	if err != nil {
		return err
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
		service.logger.Error("failed to save message to database", "error", err)
	}

	reactions := ReminderReactions
	if event.RepeatInterval != nil {
		reactions = append(reactions, ReminderReactionsRecurring...)
	}

	for _, reaction := range reactions {
		err := service.messenger.SendReactionAsync(&messenger.Reaction{
			Reaction:                  reaction,
			ChannelExternalIdentifier: room.RoomID,
			MessageExternalIdentifier: resp.ExternalIdentifier,
		})
		if err != nil {
			service.logger.Error("failed to send reaction", "matrix.reaction", reaction, "error", err)
			continue
		}
	}

	return nil
}

func (service *service) SendDailyReminder(reminder *daemon.DailyReminder, output *daemon.Output) error {
	room, err := service.matrixDatabase.GetRoomByID(output.OutputID)
	if err != nil {
		return err
	}

	msg, msgFormatted := format.InfoFromDaemonEvents(reminder.Events, room.TimeZone)
	msg = "Your Events for Today\n\n" + msg
	msgFormatted = "<h2>Your Events for Today</h2><br>\n" + msgFormatted

	resp, err := service.messenger.SendMessage(messenger.HTMLMessage(
		msg,
		msgFormatted,
		room.RoomID,
	))
	if err != nil {
		return err
	}

	dbMsg := &matrixdb.MatrixMessage{
		ID:            resp.ExternalIdentifier,
		UserID:        nil, // There can be many users in a channel.
		RoomID:        room.ID,
		Body:          msg,
		BodyFormatted: msgFormatted,
		SendAt:        resp.Timestamp,
		Type:          matrixdb.MessageTypeDailyReminder,
	}
	_, err = service.matrixDatabase.NewMessage(dbMsg)
	if err != nil {
		service.logger.Error("failed to save message to database", "error", err)
	}

	return nil
}
