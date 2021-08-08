package matrixsyncer

import (
	"fmt"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
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

// newReminder parses a message for a reminder date
func (s *Syncer) newReminder(evt *event.Event, channel *database.Channel) (*database.Reminder, error) {
	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		return nil, errors.ErrMatrixEventWrongType
	}

	remindTime, err := formater.ParseTime(content.Body, channel)
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
	if err != nil {
		log.Warn("Was not able to save a message to the database: " + err.Error())
	}

	msg := fmt.Sprintf("Successfully added new reminder (ID: %d) for %s", reminder.ID, formater.ToLocalTime(reminder.RemindTime, channel))

	log.Info(msg)
	_, err = s.messenger.SendReplyToEvent(msg, evt, evt.RoomID.String())
	if err != nil {
		log.Warn("Was not able to send success message to user")
	}

	return reminder, err
}
