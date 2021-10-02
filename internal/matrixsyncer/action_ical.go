package matrixsyncer

import (
	"strconv"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getActionIcal() *types.Action {
	action := &types.Action{
		Name:     "Get iCal link",
		Examples: []string{"ical", "calendar link", "show me the calendar link please"},
		Regex:    "(?i)(^ical$|(show|give|list|send|write|).*(calendar|ical|cal|reminder|ics)[ ]+(link|url|uri|file))",
		Action:   s.actionIcal,
	}
	return action
}

// actionList performs the action "list" that writes all pending reminders to the given channel
func (s *Syncer) actionIcal(evt *event.Event, channel *database.Channel) error {
	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		log.Warn("Event is not a message event. Can not handle it")
		return errors.ErrMatrixEventWrongType
	}

	_, err := s.daemon.Database.AddMessageFromMatrix(evt.ID.String(), time.Now().Unix(), content, nil, database.MessageTypeIcalLinkRequest, channel)
	if err != nil {
		log.Warn("Can not save message to database.")
	}

	msg := formater.Formater{}

	if len(channel.CalendarSecret) < 20 {
		msg.Text("This channel does not support calendar links. Ask your administrator to set a secret/token for you.")
	} else {
		msg.TextLine("With this link you can get access to the calendar (ics) file. Keep it secret!")
		msg.Spoiler(s.baseURL + "/calendar/" + strconv.FormatUint(uint64(channel.ID), 10) + "/ical?token=" + channel.CalendarSecret)
	}

	message, messageFormatted := msg.Build()

	_, err = s.messenger.SendFormattedMessage(message, messageFormatted, channel, database.MessageTypeIcalLink, 0)
	return err
}
