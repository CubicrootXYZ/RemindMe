package matrixsyncer

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

func (s *Syncer) getActionIcalRegenerate() *types.Action {
	action := &types.Action{
		Name:     "Renew calendar secret",
		Examples: []string{"renew the calendar secret", "generate token"},
		Regex:    "(?i)(make|generate|)[ ]*(renew|generate|delete|regenerate|renew|new)[ ]*(the|a|)[ ]+(ical|calendar|token|secret)",
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

	msg := formater.Formater{}
	err = s.daemon.Database.GenerateNewCalendarSecret(channel)
	if err != nil {
		log.Error(err.Error())
		msg.Text("Failed to generate a new secret.")
	}

	msg.Text("Updated your calendar secret. Your old secret will no longer work.")

	message, messageFormatted := msg.Build()

	_, err = s.messenger.SendFormattedMessage(message, messageFormatted, channel, database.MessageTypeIcalRenew, 0)
	return err
}
