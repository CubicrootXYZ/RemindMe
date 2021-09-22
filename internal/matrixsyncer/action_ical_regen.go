package matrixsyncer

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getActionIcalRegenerate() *Action {
	action := &Action{
		Name:     "Renew calendar secret",
		Examples: []string{"renew the calendar secret", "generate token"},
		Regex:    "(?i)(make|generate|)[ ]*(renew|generate|delete|regenerate|renew|new)[ ]*(the|a|)[ ]+(ical|calendar|token|secret)",
		Action:   s.actionIcalRegenerate,
	}
	return action
}

// actionList performs the action "list" that writes all pending reminders to the given channel
func (s *Syncer) actionIcalRegenerate(evt *event.Event, channel *database.Channel) error {
	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		log.Warn("Event is not a message event. Can not handle it")
		return errors.ErrMatrixEventWrongType
	}

	_, err := s.daemon.Database.AddMessageFromMatrix(evt.ID.String(), time.Now().Unix(), content, nil, database.MessageTypeIcalRenewRequest, channel)
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
