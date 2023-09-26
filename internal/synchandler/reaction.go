package synchandler

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// ReactionHandler handles message events
type ReactionHandler struct {
	database  types.Database
	messenger asyncmessenger.Messenger
	botInfo   *types.BotInfo
	actions   []*types.ReactionAction
	started   int64
}

// NewReactionHandler returns a new ReactionHandler
func NewReactionHandler(database types.Database, messenger asyncmessenger.Messenger, botInfo *types.BotInfo, actions []*types.ReactionAction) *ReactionHandler {
	return &ReactionHandler{
		database:  database,
		messenger: messenger,
		botInfo:   botInfo,
		actions:   actions,
		started:   time.Now().Unix(),
	}
}

// NewEvent takes a new matrix event and handles it
func (s *ReactionHandler) NewEvent(_ mautrix.EventSource, evt *event.Event) {
	log.Debug(fmt.Sprintf("New reaction: / Sender: %s / Room: / %s / Time: %d", evt.Sender, evt.RoomID, evt.Timestamp))

	// Do not answer our own and old messages
	if evt.Sender == id.UserID(s.botInfo.BotName) || evt.Timestamp/1000 < s.started {
		return
	}

	// Get all meta data
	channel, err := s.database.GetChannelByUserAndChannelIdentifier(evt.Sender.String(), evt.RoomID.String())
	if err != nil || channel == nil {
		log.Warn("Do not know that user and channel.")
	}

	content, ok := evt.Content.Parsed.(*event.ReactionEventContent)
	if !ok {
		log.Warn("Event is not a reaction event. Can not handle it.")
		return
	}

	if content.RelatesTo.EventID.String() == "" {
		log.Warn("Reaction with no realting message. Can not handle that.")
		return
	}

	message, err := s.database.GetMessageByExternalID(content.RelatesTo.EventID.String())
	if err != nil {
		log.Info("Do not know the message related to the reaction.")
		return
	}

	// Cycle through all actions
	for _, action := range s.actions {
		log.Debug("Checking for match with action " + action.Name)
		if action.Type != types.ReactionActionType(message.Type) && action.Type != types.ReactionActionTypeAll {
			continue
		}

		for _, key := range action.Keys {
			if content.RelatesTo.Key == key {
				err = action.Action(message, content, evt, channel)
				if err == nil {
					return
				}
				log.Error(err.Error())
			}
		}
	}

	log.Info("Nothing handled that reaction")
}
