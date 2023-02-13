package matrix

import (
	"errors"
	"strings"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

type MessageEvent struct {
	Event       *event.Event
	Content     *event.MessageEventContent
	IsEncrypted bool
	Room        *matrixdb.MatrixRoom
}

func (service *service) MessageEventHandler(source mautrix.EventSource, evt *event.Event) {
	logger := service.logger.WithFields(map[string]any{
		"sender":          evt.Sender,
		"room":            evt.RoomID,
		"event_timestamp": evt.Timestamp,
	})
	logger.Debugf("new message received")

	// Do not answer our own and old messages
	if evt.Sender.String() == service.botname || evt.Timestamp/1000 <= service.lastMessageFrom.Unix() {
		return
	}

	room, err := service.matrixDatabase.GetRoomByRoomID(string(evt.RoomID))
	if err != nil {
		logger.Debugf("do not know room, ignoring message")
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
		logger.Debugf("do not know user, ignoring message")
		return
	}

	// Check if we already know the message
	_, err = service.matrixDatabase.GetMessageByID(evt.ID.String())
	if err == nil {
		return
	} else {
		if !errors.Is(err, database.ErrNotFound) {
			logger.Err(err)
		}
	}

	msgEvt, err := service.parseMessageEvent(evt)
	if err != nil {
		logger.Infof("can not handle event: " + err.Error())
		return
	}
	msgEvt.Room = room

	if msgEvt.Content.RelatesTo != nil && msgEvt.Content.RelatesTo.InReplyTo != nil {
		// it is a reply
		service.findMatchingReplyAction(msgEvt, room, logger)
	} else {
		// it is a message
		service.findMatchingMessageAction(msgEvt, room, logger)
	}
}

func (service *service) findMatchingReplyAction(msgEvent *MessageEvent, room *database.MatrixRoom, logger gologger.Logger) {
	replyToMessage, err := service.matrixDatabase.GetMessageByID(msgEvent.Content.RelatesTo.InReplyTo.EventID.String())
	if err != nil {
		logger.Infof("discarding message, can not find the message it replies to: %s", err.Error())
		return
	}

	msg := strings.ToLower(msgEvent.Content.Body)
	for i := range service.config.ReplyActions {
		if service.config.ReplyActions[i].Selector().MatchString(msg) {
			logger.Infof("moving event to reply action: %s", service.config.ReplyActions[i].Name())
			service.config.ReplyActions[i].HandleEvent(msgEvent, replyToMessage)
			return
		}
	}

	logger.Infof("moving event to default reply action")
	service.config.DefaultReplyAction.HandleEvent(msgEvent, replyToMessage)
}

func (service *service) findMatchingMessageAction(msgEvent *MessageEvent, room *database.MatrixRoom, logger gologger.Logger) {
	msg := strings.ToLower(msgEvent.Content.Body)
	for i := range service.config.MessageActions {
		if service.config.MessageActions[i].Selector().MatchString(msg) {
			logger.Infof("moving event to message action: %s", service.config.MessageActions[i].Name())
			service.config.MessageActions[i].HandleEvent(msgEvent)
			return
		}
	}

	logger.Infof("moving event to default message action")
	service.config.DefaultMessageAction.HandleEvent(msgEvent)
}

func (service *service) parseMessageEvent(evt *event.Event) (*MessageEvent, error) {
	msgEvt := MessageEvent{
		Event: evt,
	}

	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if ok {
		msgEvt.Content = content
		msgEvt.IsEncrypted = false
		return &msgEvt, nil
	}

	if !service.crypto.enabled {
		return nil, ErrUnknowEvent
	}

	_, ok = evt.Content.Parsed.(*event.EncryptedEventContent)
	if ok {
		decrypted, err := service.crypto.olm.DecryptMegolmEvent(evt)

		if err != nil {
			return nil, err
		}

		content, ok = decrypted.Content.Parsed.(*event.MessageEventContent)
		if ok {
			msgEvt.Content = content
			msgEvt.IsEncrypted = true
			return &msgEvt, nil
		}
	}

	return nil, ErrUnknowEvent
}
