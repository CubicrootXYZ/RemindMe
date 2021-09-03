package matrixsyncer

import (
	"strings"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getActionTimezone() *Action {
	action := &Action{
		Name:     "Set my timezone",
		Examples: []string{"set timezone Europe/Berlin", "set timezone America/Metropolis", "set timezone Asia/Shanghai"},
		Regex:    "(?i)^set timezone .*$",
		Action:   s.actionTimezone,
	}
	return action
}

// actionList performs the action "list" that writes all pending reminders to the given channel
func (s *Syncer) actionTimezone(evt *event.Event, channel *database.Channel) error {
	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		log.Warn("Event is not a message event. Can not handle it")
		return errors.ErrMatrixEventWrongType
	}

	_, err := s.daemon.Database.AddMessageFromMatrix(evt.ID.String(), evt.Timestamp, content, nil, database.MessageTypeTimezoneChangeRequest, channel)
	if err != nil {
		log.Warn("Failed to save message in database: " + err.Error())
	}

	tz := strings.ReplaceAll(content.Body, "set timezone ", "")
	_, err = time.LoadLocation(tz)
	if err != nil {
		_, err = s.messenger.SendReplyToEvent("Sorry, I do not know this timezone.", evt, channel, database.MessageTypeTimezoneChangeRequestFail)
		return err
	}

	channel, err = s.daemon.Database.UpdateChannel(channel.ID, tz, channel.DailyReminder)
	if err != nil {
		log.Warn("Failed to save timezone in database: " + err.Error())
		_, err = s.messenger.SendReplyToEvent("Sorry, that failed.", evt, channel, database.MessageTypeTimezoneChangeRequestFail)
		return err
	}

	_, err = s.messenger.SendReplyToEvent("Great, I updated your timezone to "+tz+". Currently it is "+formater.ToLocalTime(time.Now(), channel), evt, channel, database.MessageTypeTimezoneChangeRequestSuccess)

	return err
}
