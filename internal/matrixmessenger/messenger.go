package matrixmessenger

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
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
	newMsg := fmt.Sprintf("Hey %s: %s (at %s)", "USER", reminder.Message, reminder.RemindTime.Format("15:04 02.01.2006"))
	newMsgFormatted := fmt.Sprintf("Hey %s: <br>%s <br><i>(at %s)</i>", makeLinkToUser(reminder.Channel.UserIdentifier), reminder.Message, reminder.RemindTime.Format("15:04 02.01206"))

	body, bodyFormatted := makeResponse(newMsg, newMsgFormatted, reminder.Message, reminder.Message, reminder.Channel.UserIdentifier, reminder.Channel.ChannelIdentifier, respondToMessage.ExternalIdentifier)

	matrixMessage := MatrixMessage{
		Body:          body,
		FormattedBody: bodyFormatted,
		MsgType:       "m.text",
		Type:          "m.room.message",
		Format:        "org.matrix.custom.html",
	}

	matrixMessage.RelatesTo.InReplyTo.EventID = respondToMessage.ExternalIdentifier

	evt, err := m.sendMessage(&matrixMessage, "", reminder.Channel.ChannelIdentifier)
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

// sendMessage sends a message to a matrix room
func (m *Messenger) sendMessage(message *MatrixMessage, replyEventID, roomID string) (resp *mautrix.RespSendEvent, err error) {
	log.Info(fmt.Sprintf("Sending message to room %s", roomID))
	if replyEventID != "" {
		message.RelatesTo.InReplyTo.EventID = replyEventID
	}
	return m.client.SendMessageEvent(id.RoomID(roomID), event.EventMessage, &message)
}
