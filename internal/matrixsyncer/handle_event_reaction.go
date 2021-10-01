package matrixsyncer

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func (s *Syncer) handleReactionEvent(source mautrix.EventSource, evt *event.Event) {
	log.Debug(fmt.Sprintf("New reaction: / Sender: %s / Room: / %s / Time: %d", evt.Sender, evt.RoomID, evt.Timestamp))

	// Do not answer our own and old messages
	if evt.Sender == id.UserID(s.botInfo.BotName) || evt.Timestamp/1000 <= time.Now().Unix()-60 {
		return
	}

	channel, err := s.daemon.Database.GetChannelByUserAndChannelIdentifier(evt.Sender.String(), evt.RoomID.String())
	if err != nil {
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

	message, err := s.daemon.Database.GetMessageByExternalID(content.RelatesTo.EventID.String())
	if err != nil {
		log.Info("Do not know the message related to the reaction.")
		return
	}

	for _, action := range s.reactionActions {
		log.Info("Checking for match with action " + action.Name)
		if action.Type != ReactionActionType(message.Type) && action.Type != ReactionActionTypeAll {
			continue
		}

		for _, key := range action.Keys {
			if content.RelatesTo.Key == key {
				err = action.Action(message, content, evt, channel)
				if err == nil {
					return
				}
			}
		}
	}

	log.Info("Nothing handled that reaction")
}
