package matrix

import (
	"errors"
	"log/slog"
	"strings"

	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

type MessageEvent struct {
	Event   *event.Event
	Content *event.MessageEventContent
	Room    *matrixdb.MatrixRoom
	Input   *database.Input
	Channel *database.Channel
}

func (service *service) MessageEventHandler(_ mautrix.EventSource, evt *event.Event) {
	logger := service.logger.With(
		"matrix.sender", evt.Sender,
		"matrix.room.id", evt.RoomID,
		"matrix.event.timestamp", evt.Timestamp,
	)
	logger.Debug("new message received")

	service.metricEventInCount.
		WithLabelValues("message").
		Inc()

	// Do not answer our own and old messages
	if evt.Sender.String() == service.botname || evt.Timestamp/1000 <= service.lastMessageFrom.Unix() {
		return
	}

	room, err := service.matrixDatabase.GetRoomByRoomID(string(evt.RoomID))
	if err != nil {
		logger.Debug("ignoring message", "reason", "unknown room")
		return
	}

	isUserKnown := false

	userStr := evt.Sender.String()
	for i := range room.Users {
		if room.Users[i].ID == userStr {
			isUserKnown = true
			break
		}
	}

	if !isUserKnown {
		logger.Debug("ignoring message", "reason", "unknown user")
		return
	}

	// Check if we already know the message
	msg, err := service.matrixDatabase.GetMessageByID(evt.ID.String())
	if err == nil && msg != nil {
		return
	}

	if !errors.Is(err, matrixdb.ErrNotFound) {
		logger.Error("failed to get message from database", "error", err)
	}

	msgEvt, err := service.parseMessageEvent(evt, room)
	if err != nil {
		logger.Info("failed to parse message event", "error", err)
		return
	}

	if msgEvt.Content.RelatesTo != nil && msgEvt.Content.RelatesTo.InReplyTo != nil {
		// it is a reply
		service.findMatchingReplyAction(msgEvt, logger)
	} else {
		// it is a message
		service.findMatchingMessageAction(msgEvt, logger)
	}
}

func (service *service) findMatchingReplyAction(msgEvent *MessageEvent, logger *slog.Logger) {
	replyToMessage, err := service.matrixDatabase.GetMessageByID(msgEvent.Content.RelatesTo.InReplyTo.EventID.String())
	if err != nil {
		logger.Info("failed to get replied to message from database", "error", err,
			"matrix.event.id", msgEvent.Content.RelatesTo.InReplyTo.EventID.String())

		return
	}

	msg := strings.ToLower(format.StripReply(msgEvent.Content.Body))
	for i := range service.config.ReplyActions {
		if service.config.ReplyActions[i].Selector().MatchString(msg) {
			logger.Info("moving event to reply action", "action.name", service.config.ReplyActions[i].Name())
			service.config.ReplyActions[i].HandleEvent(msgEvent, replyToMessage)

			return
		}
	}

	logger.Info("moving event to default reply action")
	service.config.DefaultReplyAction.HandleEvent(msgEvent, replyToMessage)
}

func (service *service) findMatchingMessageAction(msgEvent *MessageEvent, logger *slog.Logger) {
	msg := strings.ToLower(msgEvent.Content.Body)
	for i := range service.config.MessageActions {
		if service.config.MessageActions[i].Selector().MatchString(msg) {
			logger.Info("moving event to message action", "action.name", service.config.MessageActions[i].Name())
			service.config.MessageActions[i].HandleEvent(msgEvent)

			return
		}
	}

	logger.Info("moving event to default message action")
	service.config.DefaultMessageAction.HandleEvent(msgEvent)
}

func (service *service) parseMessageEvent(evt *event.Event, room *matrixdb.MatrixRoom) (*MessageEvent, error) {
	msgEvt := MessageEvent{
		Event: evt,
		Room:  room,
	}

	input, err := service.database.GetInputByType(room.ID, InputType)
	if err != nil {
		return nil, err
	}

	channel, err := service.database.GetChannelByID(input.ChannelID)
	if err != nil {
		return nil, err
	}

	msgEvt.Channel = channel
	msgEvt.Input = input

	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if ok {
		msgEvt.Content = content
		return &msgEvt, nil
	}

	return nil, ErrUnknowEvent
}
