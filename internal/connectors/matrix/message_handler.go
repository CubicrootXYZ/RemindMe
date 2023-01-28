package matrix

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type MessageEvent struct {
	Event       *event.Event
	Content     *event.MessageEventContent
	IsEncrypted bool
}

func (service *service) MessageEventHandler(source mautrix.EventSource, evt *event.Event) {
	logger := service.logger.WithFields(map[string]any{
		"sender":          evt.Sender,
		"room":            evt.RoomID,
		"event_timestamp": evt.Timestamp,
	})
	logger.Debugf("new message received")

	// Do not answer our own and old messages
	if evt.Sender == id.UserID(service.config.Username) { // TODO actually read last evt time from db || evt.Timestamp/1000 <= s.started {
		return
	}

	room, err := service.matrixDatabase.GetRoomByID(string(evt.RoomID))
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

	msgEvt, err := service.parseMessageEvent(evt)
	if err != nil {
		logger.Infof("can not handle event: " + err.Error())
		return
	}

	// Check if it is a reply to a message we know
	if s.checkReplyActions(msgEvt, channel) {
		return
	}

	// Check if a action matches
	if s.checkActions(msgEvt, channel) {
		return
	}

	// Nothing left so it must be a reminder
	_, err = s.newReminder(msgEvt, channel)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed parsing the Reminder with: %s", err.Error()))
		return
	}
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
