package matrix

import (
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// TODO increase test coverage.

type ReactionEvent struct {
	Event       *event.Event
	Content     *event.ReactionEventContent
	IsEncrypted bool
	Room        *matrixdb.MatrixRoom
	Input       *database.Input
	Channel     *database.Channel
}

func (service *service) ReactionEventHandler(_ mautrix.EventSource, evt *event.Event) {
	logger := service.logger.WithFields(map[string]any{
		"sender":          evt.Sender,
		"room":            evt.RoomID,
		"event_timestamp": evt.Timestamp,
	})
	logger.Debugf("new reaction received")

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

	content, ok := evt.Content.Parsed.(*event.ReactionEventContent)
	if !ok {
		logger.Infof("Event is not a reaction event. Can not handle it.")
		return
	}

	if content.RelatesTo.EventID.String() == "" {
		logger.Infof("Reaction with no relating message. Can not handle that.")
		return
	}

	message, err := service.matrixDatabase.GetMessageByID(content.RelatesTo.EventID.String())
	if err != nil {
		logger.Infof("Do not know the message related to the reaction.")
		return
	}

	if message.Room.RoomID != room.RoomID {
		// Should never happen.
		logger.Infof("Ignore reaction from room %s referring to event from room %s.", message.Room.RoomID, room.RoomID)
		return
	}

	reactionEvent, err := service.parseReactionEvent(evt, room)
	if err != nil {
		logger.Errorf("Can not parse reaction event: %s", err.Error())
		return
	}

	// Find fitting action.
	for i := range service.config.ReactionActions {
		for _, reaction := range service.config.ReactionActions[i].Selector() {
			if reaction == content.RelatesTo.Key {
				service.config.ReactionActions[i].HandleEvent(reactionEvent, message)
				return
			}
		}
	}

	logger.Infof("No action found matching key %s", content.RelatesTo.Key)
}

func (service *service) parseReactionEvent(evt *event.Event, room *matrixdb.MatrixRoom) (*ReactionEvent, error) { //nolint: dupl
	reactionEvent := ReactionEvent{
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
	reactionEvent.Channel = channel
	reactionEvent.Input = input

	content, ok := evt.Content.Parsed.(*event.ReactionEventContent)
	if ok {
		reactionEvent.Content = content
		reactionEvent.IsEncrypted = false
		return &reactionEvent, nil
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

		content, ok = decrypted.Content.Parsed.(*event.ReactionEventContent)
		if ok {
			reactionEvent.Content = content
			reactionEvent.IsEncrypted = true
			return &reactionEvent, nil
		}
	}

	return nil, ErrUnknowEvent
}
