package matrixsyncer

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/tj/go-naturaldate"
	"maunium.net/go/mautrix/event"
)

// createChannel creates a new matrix channel
func (s *Syncer) createChannel(userID string) (*database.Channel, error) {
	roomCreated, err := s.messenger.CreateChannel(userID)
	if err != nil {
		return nil, err
	}

	return s.daemon.Database.AddChannel(userID, roomCreated.RoomID.String())
}

// parseRemind parses a message for a reminder date
func (s *Syncer) parseRemind(evt *event.Event, channel *database.Channel) (*database.Reminder, error) {
	baseTime := time.Now().UTC()
	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		return nil, errors.MatrixEventWrongType
	}
	remindTime, err := naturaldate.Parse(content.Body, baseTime, naturaldate.WithDirection(naturaldate.Future))
	if err != nil {
		s.messenger.SendReplyToEvent("Sorry I was not able to understand the remind date and time from this message", evt, evt.RoomID.String())
		return nil, err
	}

	reminder, err := s.daemon.Database.AddReminder(remindTime, content.Body, true, uint64(0), channel)
	if err != nil {
		log.Warn("Error when inserting reminder: " + err.Error())
		return reminder, err
	}
	_, err = s.daemon.Database.AddMessageFromMatrix(evt.ID.String(), evt.Timestamp/1000, content, reminder, database.MessageTypeReminderRequest, channel)

	return reminder, err
}

func (s *Syncer) handleReminderRequestReaction(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) (matched bool, err error) {
	for _, action := range s.reactionActions {
		log.Info("Checking for match with action " + action.Name)
		if action.Type != ReactionActionTypeReminderRequest {
			continue
		}

		for _, key := range action.Keys {
			if content.RelatesTo.Key == key {
				err := action.Action(message, content, evt, channel)
				return true, err
			}
		}
	}
	return false, nil
}
