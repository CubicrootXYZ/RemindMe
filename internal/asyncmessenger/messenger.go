package asyncmessenger

import (
	"sync"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/encryption"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"
)

type messenger struct {
	config      *configuration.Config
	client      *mautrix.Client
	db          types.Database
	olm         *crypto.OlmMachine
	stateStore  *encryption.StateStore
	cryptoMutex sync.Mutex // Since the crypto foo relies on a single sqlite, only one process at a time is allowed to access it
	debug       bool
}

func NewMessenger(debug bool, config *configuration.Config, db types.Database, cryptoStore crypto.Store, stateStore *encryption.StateStore, matrixClient *mautrix.Client) (Messenger, error) {
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

	return &messenger{
		config:      config,
		client:      matrixClient,
		db:          db,
		olm:         olm,
		stateStore:  stateStore,
		debug:       debug,
		cryptoMutex: sync.Mutex{},
	}, nil
}

func (messenger *messenger) SendReminder(reminder *Reminder, respondToMessage *Message) (*Message, error) {
	// TODO make async? But we want to store stuff later in the database, so we need to make it async one layer above?

	// Channel is deleted do not send message
	if reminder.Channel.ID == 0 || reminder.Channel.ChannelIdentifier == "" {
		return nil, errors.ErrEmptyChannel
	}

	responseMessage, responseMessageFormatted := reminder.getRemindMessage()

	response := Response{
		Message:                   responseMessage,
		MessageFormatted:          responseMessageFormatted,
		RespondToMessage:          reminder.Message,
		RespondToMessageFormatted: reminder.Message,
		RespondToUserID:           reminder.Channel.UserIdentifier,
		RoomID:                    reminder.Channel.ChannelIdentifier,
		RespondToEventID:          respondToMessage.ExternalIdentifier,
	}

	matrixMessage := response.toMatrixMessage()

	// TODO continue

	evt, err := m.sendMessage(&matrixMessage, reminder.Channel.ChannelIdentifier, event.EventMessage)
	if err != nil {
		return nil, err
	}

	for _, reaction := range types.ReactionsReminder {
		_, err = m.SendReaction(reaction, string(evt.EventID), &reminder.Channel)
		time.Sleep(time.Millisecond * 500) // Avoid getting blocked
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
