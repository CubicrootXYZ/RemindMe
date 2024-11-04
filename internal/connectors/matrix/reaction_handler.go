package matrix

import (
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// TODO increase test coverage.

type ReactionEvent struct {
	Event   *event.Event
	Content *event.ReactionEventContent
	Room    *matrixdb.MatrixRoom
	Input   *database.Input
	Channel *database.Channel
}

func (service *service) ReactionEventHandler(_ mautrix.EventSource, evt *event.Event) {
	logger := service.logger.With(
		"matrix.sender", evt.Sender,
		"matrix.room.id", evt.RoomID,
		"matrix.event.timestamp", evt.Timestamp,
	)
	logger.Debug("new reaction received")

	// Do not answer our own and old messages
	if evt.Sender.String() == service.botname || evt.Timestamp/1000 <= service.lastMessageFrom.Unix() {
		return
	}

	room, err := service.matrixDatabase.GetRoomByRoomID(string(evt.RoomID))
	if err != nil {
		logger.Debug("ignoring reaction", "reason", "unknown room")
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
		logger.Debug("ignoring reaction", "reason", "unknown user")
		return
	}

	content, ok := evt.Content.Parsed.(*event.ReactionEventContent)
	if !ok {
		logger.Info("ignoring reaction", "reason", "not a reaction event")
		return
	}

	if content.RelatesTo.EventID.String() == "" {
		logger.Info("ignoring reaction", "reason", "no related event")
		return
	}

	message, err := service.matrixDatabase.GetMessageByID(content.RelatesTo.EventID.String())
	if err != nil {
		logger.Info("ignoring reaction", "reason", "unknown related message")
		return
	}

	if message.Room.RoomID != room.RoomID {
		// Should never happen.
		logger.Info("ignoring reaction", "reason", "message from different room than reaction",
			"matrix.message.room.id", message.Room.RoomID, "matrix.reaction.room.id", room.RoomID)
		return
	}

	reactionEvent, err := service.parseReactionEvent(evt, room)
	if err != nil {
		logger.Error("failed to parse reaction event", "error", err)
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

	logger.Info("ignoring reaction", "reason", "unknown reaction", "matrix.reaction.key", content.RelatesTo.Key)
}

func (service *service) parseReactionEvent(evt *event.Event, room *matrixdb.MatrixRoom) (*ReactionEvent, error) {
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
		return &reactionEvent, nil
	}

	return nil, ErrUnknowEvent
}
