package matrixsyncer

import (
	"regexp"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

func (s *Syncer) getActionIcalRegenerate() *types.Action {
	action := &types.Action{
		Name:     "Renew calendar secret",
		Examples: []string{"renew the calendar secret", "generate token"},
		Regex:    regexp.MustCompile("(?i)^(make|generate|)[ ]*(renew|generate|delete|regenerate|renew|new)[ ]*(the|a|)[ ]+(ical|calendar|token|secret)[ ]*(token|secret|)[ ]*$"),
		Action:   s.actionIcalRegenerate,
	}
	return action
}

// actionList performs the action "list" that writes all pending reminders to the given channel
func (s *Syncer) actionIcalRegenerate(evt *types.MessageEvent, channel *database.Channel) error {
	_, err := s.daemon.Database.AddMessageFromMatrix(evt.Event.ID.String(), time.Now().Unix(), evt.Content, nil, database.MessageTypeIcalRenewRequest, channel)
	if err != nil {
		log.Warn("Can not save message to database.")
	}

	msg := "Updated your calendar secret. Your old secret will no longer work."
	err = s.daemon.Database.GenerateNewCalendarSecret(channel)
	if err != nil {
		log.Error(err.Error())
		msg = "Failed to generate a new secret."
	}

	go s.sendAndStoreMessage(asyncmessenger.PlainTextMessage(
		msg,
		channel.ChannelIdentifier,
	), channel, database.MessageTypeIcalRenew, 0)
	return nil
}
