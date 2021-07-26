package matrixmessenger

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/dchest/uniuri"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Messenger holds all information for messaging
type Messenger struct {
	config *configuration.Matrix
	client *mautrix.Client
}

// MatrixMessage holds information for a matrix response message
type MatrixMessage struct {
	Body          string `json:"body"`
	Format        string `json:"format"`
	FormattedBody string `json:"formatted_body,omitempty"`
	MsgType       string `json:"msgtype"`
	Type          string `json:"type"`
	RelatesTo     struct {
		InReplyTo struct {
			EventID string `json:"event_id,omitempty"`
		} `json:"m.in_reply_to,omitempty"`
	} `json:"m.relates_to,omitempty"`
}

// Create creates a new matrix messenger
func Create(config *configuration.Matrix) (*Messenger, error) {
	// Log into matrix
	client, err := mautrix.NewClient(config.Homeserver, "", "")
	if err != nil {
		return nil, err
	}

	_, err = client.Login(&mautrix.ReqLogin{
		Type:             "m.login.password",
		Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: config.Username},
		Password:         config.Password,
		StoreCredentials: true,
	})
	if err != nil {
		return nil, err
	}
	log.Info("Logged in to matrix")

	return &Messenger{
		config: config,
		client: client,
	}, nil
}

// SendReminder sends a reminder to the user
func (m *Messenger) SendReminder(reminder *database.Reminder, respondToMessage *database.Message) (*database.Message, error) {
	newMsg := fmt.Sprintf("%s a reminder for you: %s (at %s)", "USER", reminder.Message, reminder.RemindTime.Format("15:04 02.01.2006"))
	newMsgFormatted := fmt.Sprintf("%s a Reminder for you: <br>%s <br><i>(at %s)</i>", makeLinkToUser(reminder.Channel.UserIdentifier), reminder.Message, reminder.RemindTime.Format("15:04 02.01.2006"))

	body, bodyFormatted := makeResponse(newMsg, newMsgFormatted, reminder.Message, reminder.Message, reminder.Channel.UserIdentifier, reminder.Channel.ChannelIdentifier, respondToMessage.ExternalIdentifier)

	matrixMessage := MatrixMessage{
		Body:          body,
		FormattedBody: bodyFormatted,
		MsgType:       "m.text",
		Type:          "m.room.message",
		Format:        "org.matrix.custom.html",
	}

	matrixMessage.RelatesTo.InReplyTo.EventID = respondToMessage.ExternalIdentifier

	evt, err := m.sendMessage(&matrixMessage, reminder.Channel.ChannelIdentifier)
	if err != nil {
		return nil, err
	}

	message := database.Message{
		Body:               matrixMessage.Body,
		BodyHTML:           matrixMessage.FormattedBody,
		ReminderID:         reminder.ID,
		Reminder:           *reminder,
		Type:               database.MessageTypeReminder,
		ChannelID:          reminder.ChannelID,
		Channel:            reminder.Channel,
		Timestamp:          time.Now().Unix(),
		ExternalIdentifier: evt.EventID.String(),
	}
	return &message, nil
}

// SendReplyToEvent sends a message in reply to the given replyEvent, if the event is nil or of wrogn format a normal message will be sent
func (m *Messenger) SendReplyToEvent(msg string, replyEvent *event.Event, roomID string) (resp *mautrix.RespSendEvent, err error) {
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

		body, bodyFormatted := makeResponse(msg, msg, content.Body, oldFormattedBody, replyEvent.Sender.String(), roomID, replyEvent.ID.String())

		message.Body = body
		message.FormattedBody = bodyFormatted
		message.RelatesTo.InReplyTo.EventID = replyEvent.ID.String()
	} else {
		message.Body = msg
		message.FormattedBody = msg
	}

	message.Format = "org.matrix.custom.html"
	message.MsgType = "m.text"
	message.Type = "m.room.message"

	return m.sendMessage(&message, roomID)
}

// CreateChannel creates a new matrix channel
func (m *Messenger) CreateChannel(userID string) (*mautrix.RespCreateRoom, error) {
	// TODO use another alias name that is more unique
	room := mautrix.ReqCreateRoom{
		Visibility:    "private",
		RoomAliasName: "RemindMe-" + uniuri.NewLen(5),
		Name:          "RemindMe",
		Topic:         "I will be your personal reminder bot",
		Invite:        []id.UserID{id.UserID(userID)},
		Preset:        "trusted_private_chat",
	}

	return m.client.CreateRoom(&room)
}

// sendMessage sends a message to a matrix room
func (m *Messenger) sendMessage(message *MatrixMessage, roomID string) (resp *mautrix.RespSendEvent, err error) {
	log.Info(fmt.Sprintf("Sending message to room %s", roomID))
	return m.client.SendMessageEvent(id.RoomID(roomID), event.EventMessage, &message)
}

// SendFormattedMessage sends a HTML formatted message to the given room
func (m *Messenger) SendFormattedMessage(msg, msgFormatted, roomID string) (resp *mautrix.RespSendEvent, err error) {
	message := MatrixMessage{
		Body:          msg,
		FormattedBody: msgFormatted,
		MsgType:       "m.text",
		Type:          "m.room.message",
		Format:        "org.matrix.custom.html",
	}

	return m.sendMessage(&message, roomID)
}
