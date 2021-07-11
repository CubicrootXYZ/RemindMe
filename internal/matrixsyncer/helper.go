package matrixsyncer

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/dchest/uniuri"
	"github.com/tj/go-naturaldate"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// createChannel creates a new matrix channel
func (s *Syncer) createChannel(userID string) (*database.Channel, error) {
	if s.client == nil {
		log.Error("Can not create a channel without a matrix client")
		return nil, errors.MatrixClientNotInitialized
	}

	// TODO use another alias name that is more unique
	room := mautrix.ReqCreateRoom{
		Visibility:    "private",
		RoomAliasName: "RemindMe-" + uniuri.NewLen(5),
		Name:          "RemindMe",
		Topic:         "I will be your personal reminder bot",
		Invite:        []id.UserID{id.UserID(userID)},
		Preset:        "trusted_private_chat",
	}

	roomCreated, err := s.client.CreateRoom(&room)
	if err != nil {
		return nil, err
	}

	return s.daemon.Database.AddChannel(userID, roomCreated.RoomID.String())
}

// sendMessage sends a message to a matrix room
func (s *Syncer) sendMessage(msg string, replyEvent *event.Event, roomID string) (resp *mautrix.RespSendEvent, err error) {
	var message MatrixMessage
	if replyEvent != nil {
		content, ok := replyEvent.Content.Parsed.(*event.MessageEventContent)
		if !ok {
			return nil, errors.MatrixEventWrongType
		}

		oldFormattedBody := content.Body
		if len(content.FormattedBody) > 1 {
			oldFormattedBody = content.FormattedBody
		}

		message.Body = fmt.Sprintf("> <%s>%s\n\n%s", replyEvent.Sender.String(), content.Body, msg)
		message.FormattedBody = fmt.Sprintf("<mx-reply><blockquote><a href='https://matrix.to/#/%s/%s'>In reply to</a> <a href='https://matrix.to/#/%s'>%s</a><br />%s</blockquote>\n</mx-reply>%s", roomID, replyEvent.ID.String(), replyEvent.Sender.String(), replyEvent.Sender.String(), oldFormattedBody, msg)
		message.Relatesto.InReplyTo.EventID = replyEvent.ID.String()
	} else {
		message.Body = msg
		message.FormattedBody = msg
	}

	message.Format = "org.matrix.custom.html"
	message.MsgType = "m.text"
	message.Type = "m.room.message"

	return s.client.SendMessageEvent(id.RoomID(roomID), event.EventMessage, &message)
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
		s.sendMessage("Sorry I was not able to understand the remind date and time from this message", evt, evt.RoomID.String())
		return nil, err
	}

	reminder, err := s.daemon.Database.AddReminder(remindTime, content.Body, true, uint64(0), channel)
	if err != nil {
		log.Warn("Error when inserting reminder: " + err.Error())
		return reminder, err
	}
	_, err = s.daemon.Database.AddMessage(evt.ID.String(), evt.Timestamp/1000, content, reminder, database.MessageTypeReminderRequest, channel)

	return reminder, err
}
