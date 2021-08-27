package matrixsyncer

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/gonaturalduration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getReplyActionRecurring(rtt database.MessageType) *ReplyAction {
	action := &ReplyAction{
		Name:        "Make a reminder recurring",
		Examples:    []string{"every 10 days", "each twenty two hours and five seconds"},
		Regex:       "(?i)(every|each|always|recurring|all|any).*(second|minute|day|hour)(|s)$",
		ReplyToType: rtt,
		Action:      s.replyActionRecurring,
	}
	return action
}

func (s *Syncer) replyActionRecurring(evt *event.Event, channel *database.Channel, replyMessage *database.Message, content *event.MessageEventContent) error {
	// Get duration from message
	duration := gonaturalduration.ParseNumber(content.Body)
	if duration <= time.Minute {
		log.Info("Duration was < 0")
		return nil
	}

	// Repeat for 5 years
	repeatTimes := (5 * 365 * 24 * time.Hour) / duration

	reminder, err := s.daemon.Database.UpdateReminder(replyMessage.ReminderID, time.Now(), uint64(duration/time.Minute), uint64(repeatTimes))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	_, err = s.daemon.Database.AddMessageFromMatrix(evt.ID.String(), evt.Timestamp, content, reminder, database.MessageTypeReminderRecurringRequest, channel)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to add recurring message %s to database: %s", evt.ID.String(), err.Error()))
	}

	lastRemind := reminder.RemindTime.Add(duration * repeatTimes)

	msg := fmt.Sprintf("Updated the reminder to remind you every %s until %s", formater.ToNiceDuration(duration), formater.ToLocalTime(lastRemind, channel))
	resp, err := s.messenger.SendReplyToEvent(msg, evt, channel.ChannelIdentifier)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	_, err = s.daemon.Database.AddMessageFromMatrix(resp.EventID.String(), time.Now().Unix(), nil, reminder, database.MessageTypeReminderRecurringSuccess, channel)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to add recurring message %s to database: %s", evt.ID.String(), err.Error()))
	}
	return err
}
