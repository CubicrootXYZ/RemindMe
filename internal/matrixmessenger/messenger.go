package matrixmessenger

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/encryption"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"github.com/dchest/uniuri"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Messenger holds all information for messaging
type Messenger struct {
	config     *configuration.Config
	client     *mautrix.Client
	db         types.Database
	olm        *crypto.OlmMachine
	stateStore *encryption.StateStore
	debug      bool
}

// MatrixMessage holds information for a matrix response message
type MatrixMessage struct {
	Body          string `json:"body,omitempty"`
	Format        string `json:"format,omitempty"`
	FormattedBody string `json:"formatted_body,omitempty"`
	MsgType       string `json:"msgtype,omitempty"`
	Type          string `json:"type"`
	RelatesTo     struct {
		EventID   string `json:"event_id,omitempty"`
		Key       string `json:"key,omitempty"`
		RelType   string `json:"rel_type,omitempty"`
		InReplyTo struct {
			EventID string `json:"event_id,omitempty"`
		} `json:"m.in_reply_to,omitempty"`
	} `json:"m.relates_to,omitempty"`
}

// Create creates a new matrix messenger
func Create(debug bool, config *configuration.Config, db types.Database, cryptoStore crypto.Store, stateStore *encryption.StateStore, matrixClient *mautrix.Client) (*Messenger, error) {
	var olm *crypto.OlmMachine
	if config.MatrixBotAccount.E2EE {
		olm = encryption.GetOlmMachine(debug, matrixClient, cryptoStore, db, stateStore)
		olm.AllowUnverifiedDevices = true
		olm.ShareKeysToUnverifiedDevices = true
		err := olm.Load()
		if err != nil {
			return nil, err
		}
	}

	return &Messenger{
		config:     config,
		client:     matrixClient,
		db:         db,
		olm:        olm,
		stateStore: stateStore,
		debug:      debug,
	}, nil
}

// SendReminder sends a reminder to the user
func (m *Messenger) SendReminder(reminder *database.Reminder, respondToMessage *database.Message) (*database.Message, error) {
	// Channel is deleted do not send message
	if reminder.Channel.ID == 0 || reminder.Channel.ChannelIdentifier == "" {
		return nil, errors.ErrEmptyChannel
	}

	newMsg := fmt.Sprintf("%s a reminder for you: %s (at %s)", "USER", reminder.Message, formater.ToLocalTime(reminder.RemindTime, &reminder.Channel))
	newMsgFormatted := fmt.Sprintf("%s a Reminder for you: <br>%s <br><i>(at %s)</i>", makeLinkToUser(reminder.Channel.UserIdentifier), reminder.Message, formater.ToLocalTime(reminder.RemindTime, &reminder.Channel))

	body, bodyFormatted := makeResponse(newMsg, newMsgFormatted, reminder.Message, reminder.Message, reminder.Channel.UserIdentifier, reminder.Channel.ChannelIdentifier, respondToMessage.ExternalIdentifier)

	matrixMessage := MatrixMessage{
		Body:          body,
		FormattedBody: bodyFormatted,
		MsgType:       "m.text",
		Type:          "m.room.message",
		Format:        "org.matrix.custom.html",
	}

	matrixMessage.RelatesTo.InReplyTo.EventID = respondToMessage.ExternalIdentifier

	evt, err := m.sendMessage(&matrixMessage, reminder.Channel.ChannelIdentifier, event.EventMessage)
	if err != nil {
		return nil, err
	}

	for _, reaction := range types.ReactionsReminderRequest {
		_, err = m.SendReaction(reaction, string(evt.EventID), &reminder.Channel)
		if err != nil {
			log.Warn(err.Error())
		}
	}

	message := database.Message{
		Body:               matrixMessage.Body,
		BodyHTML:           matrixMessage.FormattedBody,
		ReminderID:         &reminder.ID,
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
func (m *Messenger) SendReplyToEvent(msg string, replyEvent *types.MessageEvent, channel *database.Channel, msgType database.MessageType) (resp *mautrix.RespSendEvent, err error) {
	var message MatrixMessage
	if replyEvent != nil {
		oldFormattedBody := formater.StripReply(replyEvent.Content.Body)
		if len(replyEvent.Content.FormattedBody) > 1 {
			oldFormattedBody = formater.StripReplyFormatted(replyEvent.Content.FormattedBody)
		}

		body, bodyFormatted := makeResponse(msg, msg, formater.StripReply(replyEvent.Content.Body), oldFormattedBody, replyEvent.Event.Sender.String(), channel.ChannelIdentifier, replyEvent.Event.ID.String())

		message.Body = body
		message.FormattedBody = bodyFormatted
		message.RelatesTo.InReplyTo.EventID = replyEvent.Event.ID.String()
	} else {
		message.Body = msg
		message.FormattedBody = msg
	}

	message.Format = "org.matrix.custom.html"
	message.MsgType = "m.text"
	message.Type = "m.room.message"

	resp, err = m.sendMessage(&message, channel.ChannelIdentifier, event.EventMessage)
	if err != nil {
		log.Warn("Could not send message: " + err.Error())
		return nil, err
	}

	// Add message to the database
	if msgType != database.MessageTypeDoNotSave {
		dbMessage := &database.Message{
			Body:               msg,
			BodyHTML:           msg,
			ResponseToMessage:  replyEvent.Event.ID.String(),
			Type:               msgType,
			ChannelID:          channel.ID,
			Channel:            *channel,
			Timestamp:          time.Now().Unix(),
			ExternalIdentifier: resp.EventID.String(),
		}

		origMessage, err := m.db.GetMessageByExternalID(replyEvent.Event.ID.String())
		if err == nil && origMessage.ReminderID != nil && *origMessage.ReminderID != 0 {
			dbMessage.ReminderID = origMessage.ReminderID
		}

		_, err = m.db.AddMessage(dbMessage)
		if err != nil {
			log.Warn(fmt.Sprintf("Failed to save message of type %s in database: %s", string(msgType), err.Error()))
		}
	}

	return resp, err
}

// CreateChannel creates a new matrix channel
func (m *Messenger) CreateChannel(userID string) (*mautrix.RespCreateRoom, error) {
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
func (m *Messenger) sendMessage(message *MatrixMessage, roomID string, eventType event.Type) (resp *mautrix.RespSendEvent, err error) {
	if m.stateStore != nil && eventType != event.EventReaction {
		if m.stateStore.IsEncrypted(id.RoomID(roomID)) && m.olm != nil {
			resp, err = m.sendEncryptedMessage(message, roomID, eventType)
			if err == nil {
				return resp, nil
			}
		}
	}

	log.Info(fmt.Sprintf("Sending message to room %s", roomID))
	return m.client.SendMessageEvent(id.RoomID(roomID), eventType, &message)
}

func (m *Messenger) sendEncryptedMessage(message *MatrixMessage, roomID string, eventType event.Type) (resp *mautrix.RespSendEvent, err error) {
	encrypted, err := m.olm.EncryptMegolmEvent(id.RoomID(roomID), eventType, message)

	if err == crypto.SessionExpired || err == crypto.SessionNotShared || err == crypto.NoGroupSession {
		err = m.olm.ShareGroupSession(id.RoomID(roomID), m.getUserIDs(id.RoomID(roomID)))
		if err != nil {
			return nil, err
		}

		encrypted, err = m.olm.EncryptMegolmEvent(id.RoomID(roomID), eventType, message)
	}
	if err != nil {
		return nil, err
	}

	log.Info(fmt.Sprintf("Sending encrypted message to room %s", roomID))
	return m.client.SendMessageEvent(id.RoomID(roomID), event.EventEncrypted, encrypted)
}

// SendFormattedMessage sends a HTML formatted message to the given room
func (m *Messenger) SendFormattedMessage(msg, msgFormatted string, channel *database.Channel, msgType database.MessageType, relatedReminderID uint) (resp *mautrix.RespSendEvent, err error) {
	message := MatrixMessage{
		Body:          msg,
		FormattedBody: msgFormatted,
		MsgType:       "m.text",
		Type:          "m.room.message",
		Format:        "org.matrix.custom.html",
	}

	resp, err = m.sendMessage(&message, channel.ChannelIdentifier, event.EventMessage)

	// Add message to the database
	if msgType != database.MessageTypeDoNotSave {
		dbMessage := &database.Message{
			Body:               msg,
			BodyHTML:           msg,
			Type:               msgType,
			ChannelID:          channel.ID,
			Channel:            *channel,
			Timestamp:          time.Now().Unix(),
			ExternalIdentifier: resp.EventID.String(),
		}

		if relatedReminderID != 0 {
			dbMessage.ReminderID = &relatedReminderID
		}

		_, err = m.db.AddMessage(dbMessage)
		if err != nil {
			log.Warn(fmt.Sprintf("Failed to save message of type %s in database: %s", string(msgType), err.Error()))
		}
	}

	return resp, err
}

// DeleteMessage deletes a message in matrix
func (m *Messenger) DeleteMessage(messageID, roomID string) error {
	_, err := m.client.RedactEvent(id.RoomID(roomID), id.EventID(messageID))
	return err
}

// SendNotice sends a notice to the room
func (m *Messenger) SendNotice(msg, roomID string) (resp *mautrix.RespSendEvent, err error) {
	message := MatrixMessage{
		Body:    msg,
		MsgType: "m.notice",
		Type:    "m.room.message",
	}

	return m.sendMessage(&message, roomID, event.EventMessage)
}

func (m *Messenger) getUserIDs(roomID id.RoomID) []id.UserID {
	userIDs := make([]id.UserID, 0)
	members, err := m.client.JoinedMembers(roomID)
	if err != nil {
		log.Warn(err.Error())
		return userIDs
	}

	i := 0
	for userID := range members.Joined {
		userIDs = append(userIDs, userID)
		i++
	}
	return userIDs
}

func (m *Messenger) SendReaction(reaction string, toMessage string, channel *database.Channel) (resp *mautrix.RespSendEvent, err error) {
	if !m.config.BotSettings.SendReactions {
		return nil, errors.ErrReactionsDisabled
	}
	var message MatrixMessage
	message.Type = "m.reaction"
	message.RelatesTo.EventID = toMessage
	message.RelatesTo.Key = reaction
	message.RelatesTo.RelType = "m.annotation"

	return m.sendMessage(&message, channel.ChannelIdentifier, event.EventReaction)
}
