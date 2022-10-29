package matrixsyncer

import (
	"regexp"
	"strings"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

func (s *Syncer) getActionTimezone() *types.Action {
	action := &types.Action{
		Name:     "Set my timezone",
		Examples: []string{"set timezone Europe/Berlin", "set timezone America/Metropolis", "set timezone Asia/Shanghai"},
		Regex:    regexp.MustCompile("(?i)^set timezone .*$"),
		Action:   s.actionTimezone,
	}
	return action
}

// actionList performs the action "list" that writes all pending reminders to the given channel
func (s *Syncer) actionTimezone(evt *types.MessageEvent, channel *database.Channel) error {
	_, err := s.daemon.Database.AddMessageFromMatrix(evt.Event.ID.String(), evt.Event.Timestamp, evt.Content, nil, database.MessageTypeTimezoneChangeRequest, channel)
	if err != nil {
		log.Warn("Failed to save message in database: " + err.Error())
	}

	tz := strings.ReplaceAll(evt.Content.Body, "set timezone ", "")
	_, err = time.LoadLocation(tz)
	if err != nil {
		go s.sendAndStoreReply(asyncmessenger.PlainTextResponse(
			"Sorry, I do not know this timezone.",
			evt.Event.ID.String(),
			evt.Content.Body,
			evt.Event.Sender.String(),
			channel.ChannelIdentifier,
		), channel, database.MessageTypeTimezoneChangeRequestFail, 0)
		return err
	}

	channel, err = s.daemon.Database.UpdateChannel(channel.ID, tz, channel.DailyReminder, channel.Role)
	if err != nil {
		log.Warn("Failed to save timezone in database: " + err.Error())
		go s.sendAndStoreReply(asyncmessenger.PlainTextResponse(
			"Sorry, that failed.",
			evt.Event.ID.String(),
			evt.Content.Body,
			evt.Event.Sender.String(),
			channel.ChannelIdentifier,
		), channel, database.MessageTypeTimezoneChangeRequestFail, 0)
		return err
	}

	go s.sendAndStoreReply(asyncmessenger.PlainTextResponse(
		"Great, I updated your timezone to "+tz+". Currently it is "+formater.ToLocalTime(time.Now(), channel.TimeZone),
		evt.Event.ID.String(),
		evt.Content.Body,
		evt.Event.Sender.String(),
		channel.ChannelIdentifier,
	), channel, database.MessageTypeTimezoneChangeRequestSuccess, 0)

	return nil
}
