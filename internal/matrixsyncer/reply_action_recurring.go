package matrixsyncer

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/gonaturalduration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

func (s *Syncer) getReplyActionRecurring(rtt []database.MessageType) *types.ReplyAction {
	action := &types.ReplyAction{
		Name:         "Make a reminder recurring",
		Examples:     []string{"every 10 days", "each twenty two hours and five seconds"},
		Regex:        "(?i)(repeat|every|each|always|recurring|all|any).*(second|minute|hour|day|week|month|year)(|s)$",
		ReplyToTypes: rtt,
		Action:       s.replyActionRecurring,
	}
	return action
}

func (s *Syncer) replyActionRecurring(evt *types.MessageEvent, channel *database.Channel, replyMessage *database.Message) error {
	if replyMessage.ReminderID == nil {
		msg := fmt.Sprintf("Sorry, I could not delete the reminder %d.", replyMessage.ReminderID)
		msgFormatted := msg
		_, _ = s.messenger.SendFormattedMessage(msg, msgFormatted, channel, database.MessageTypeReminderRecurringFail, 0)
		return errors.ErrIDNotSet
	}

	// Get duration from message
	duration := gonaturalduration.ParseNumber(evt.Content.Body)
	if duration <= time.Minute {
		log.Info("Duration was < 0")
		return nil
	}

	// Repeat for 5 years
	repeatTimes := (5 * 365 * 24 * time.Hour) / duration

	reminder, err := s.daemon.Database.UpdateReminder(*replyMessage.ReminderID, time.Now(), uint64(duration/time.Minute), uint64(repeatTimes))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	_, err = s.daemon.Database.AddMessageFromMatrix(evt.Event.ID.String(), evt.Event.Timestamp, evt.Content, reminder, database.MessageTypeReminderRecurringRequest, channel)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to add recurring message %s to database: %s", evt.Event.ID.String(), err.Error()))
	}

	lastRemind := reminder.RemindTime.Add(duration * repeatTimes)

	msg := fmt.Sprintf("Updated the reminder to remind you every %s until %s", formater.ToNiceDuration(duration), formater.ToLocalTime(lastRemind, channel))
	_, err = s.messenger.SendReplyToEvent(msg, evt, channel, database.MessageTypeReminderRecurringSuccess)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return err
}
