package matrixsyncer

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/encryption"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/eventdaemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/synchandler"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

// Errors the service exposes to be checked by third parties
var (
	ErrMessageHasNoReminder = errors.New("need the message to be assigned to a reminder but is not")
)

// Syncer receives messages from a matrix channel
type Syncer struct {
	config      configuration.Matrix
	baseURL     string
	adminUsers  []string
	client      *mautrix.Client
	daemon      *eventdaemon.Daemon
	botSettings *configuration.BotSettings
	botInfo     *types.BotInfo
	messenger   asyncmessenger.Messenger
	cryptoStore crypto.Store
	stateStore  *encryption.StateStore
	debug       bool
}

// Create creates a new syncer
func Create(config *configuration.Config, matrixAdminUsers []string, messenger asyncmessenger.Messenger, cryptoStore crypto.Store, stateStore *encryption.StateStore, matrixClient *mautrix.Client) *Syncer {
	syncer := &Syncer{
		config:      config.MatrixBotAccount,
		baseURL:     config.Webserver.BaseURL,
		adminUsers:  matrixAdminUsers,
		messenger:   messenger,
		botSettings: &config.BotSettings,
		cryptoStore: cryptoStore,
		stateStore:  stateStore,
		debug:       config.Debug,
		client:      matrixClient,
	}

	return syncer
}

// Start starts the syncer
func (s *Syncer) Start(daemon *eventdaemon.Daemon) error {
	log.Info(fmt.Sprintf("Starting Matrix syncer for user %s on server %s", s.config.Username, s.config.Homeserver))

	s.daemon = daemon
	s.botInfo = &types.BotInfo{
		BotName: fmt.Sprintf("@%s:%s", s.config.Username, strings.ReplaceAll(strings.ReplaceAll(s.config.Homeserver, "https://", ""), "http://", "")),
	}

	s.client.Store = mautrix.NewInMemoryStore()

	var olm *crypto.OlmMachine
	if s.config.E2EE {
		olm = encryption.GetOlmMachine(s.debug, s.client, s.cryptoStore, s.daemon.Database, s.stateStore)
		olm.AllowUnverifiedDevices = true
		olm.ShareKeysToUnverifiedDevices = true
		if err := olm.Load(); err != nil {
			return err
		}
	}

	if err := s.syncChannels(); err != nil {
		return err
	}

	// Load actions
	messageActions := s.getActions()
	replyActions := s.getReplyActions()
	reactionActions := s.getReactionActions()

	// Initialize handler
	messageHandler := synchandler.NewMessageHandler(s.daemon.Database, s.messenger, s.botInfo, replyActions, messageActions, olm)
	stateMemberHandler := synchandler.NewStateMemberHandler(s.daemon.Database, s.messenger, s.client, s.botInfo, s.botSettings, olm)
	reactionHandler := synchandler.NewReactionHandler(s.daemon.Database, s.messenger, s.botInfo, reactionActions)

	// Get messages
	syncer, ok := s.client.Syncer.(*mautrix.DefaultSyncer)
	if !ok {
		return errors.New("Syncer of wrong type")
	}

	if s.config.E2EE {
		log.Info("Listening for E2EE events.")
		syncer.OnSync(func(resp *mautrix.RespSync, since string) bool {
			olm.ProcessSyncResponse(resp, since)
			return true
		})
		syncer.OnEventType(event.EventEncrypted, messageHandler.NewEvent)
		syncer.OnEventType(event.StateEncryption, func(_ mautrix.EventSource, event *event.Event) {
			s.stateStore.SetEncryptionEvent(event)
		})
	}

	syncer.OnEventType(event.EventMessage, messageHandler.NewEvent)
	syncer.OnEventType(event.EventReaction, reactionHandler.NewEvent)
	syncer.OnEventType(event.StateMember, stateMemberHandler.NewEvent)

	return s.client.Sync()
}

// Stop stops the syncer
func (s *Syncer) Stop() {
	s.client.StopSync()
}

// syncChannel keeps the channelös in sync with the config file and up to date
func (s *Syncer) syncChannels() error {
	log.Info("Syncing channels")

	channels := make([]*database.Channel, 0)

	for _, user := range s.adminUsers {
		// Get or create channel
		channel, err := s.daemon.Database.GetChannelByUserIdentifier(user)
		if err == gorm.ErrRecordNotFound {
			channel, err = s.createChannel(user, roles.RoleAdmin)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// Upgrade to newest feature set
		channel, err = s.upgradeChannel(channel, roles.RoleAdmin)
		if err != nil {
			log.Error(err.Error())
		}

		channels = append(channels, channel)

		_ = s.messenger.SendMessageAsync(asyncmessenger.PlainTextMessage(
			"Sorry I was sleeping for a while. I am now ready for your requests!",
			channel.ChannelIdentifier,
		))
	}

	// Remove channels not needed anymore
	err := s.daemon.Database.CleanAdminChannels(channels)
	if err != nil {
		log.Warn("Can not clean channels list")
		panic(err)
	}

	return nil
}

// upgradeChannel upgrades a channel to the newest features
func (s *Syncer) upgradeChannel(channel *database.Channel, defaultRole roles.Role) (*database.Channel, error) {
	var err error

	// Update secret if not set
	if len(channel.CalendarSecret) < 20 {
		err = s.daemon.Database.GenerateNewCalendarSecret(channel)
		if err != nil {
			return nil, err
		}
	}

	if channel.Role == nil {
		channel, err = s.daemon.Database.UpdateChannel(channel.ID, channel.TimeZone, channel.DailyReminder, &defaultRole)
	}

	return channel, err
}

func (s *Syncer) getActions() []*types.Action {
	messageActions := make([]*types.Action, 0)
	messageActions = append(messageActions, s.getActionList())
	messageActions = append(messageActions, s.getActionCommands())
	messageActions = append(messageActions, s.getActionTimezone())
	messageActions = append(messageActions, s.getActionSetDailyReminder())
	messageActions = append(messageActions, s.getActionDeleteDailyReminder())
	messageActions = append(messageActions, s.getActionIcal())
	messageActions = append(messageActions, s.getActionIcalRegenerate())
	messageActions = append(messageActions, s.getActionDelete())
	messageActions = append(messageActions, s.getActionAddUser())
	messageActions = append(messageActions, s.getActionDeleteReminder())
	messageActions = append(messageActions, s.getActionChangeReminder())
	return messageActions
}

func (s *Syncer) getReplyActions() []*types.ReplyAction {
	replyActions := make([]*types.ReplyAction, 0)
	replyActions = append(replyActions, s.getReplyActionDelete(database.MessageTypesWithReminder))
	replyActions = append(replyActions, s.getReplyActionRecurring(database.MessageTypesWithReminder))

	return replyActions
}

func (s *Syncer) getReactionActions() []*types.ReactionAction {
	reactionActions := make([]*types.ReactionAction, 0)
	reactionActions = append(reactionActions, s.getReactionActionDelete(types.ReactionActionTypeReminderRequest))
	reactionActions = append(reactionActions, s.getReactionActionDelete(types.ReactionActionTypeReminder))
	reactionActions = append(reactionActions, s.getReactionsAddTime(types.ReactionActionTypeReminderRequest)...)
	reactionActions = append(reactionActions, s.getReactionsAddTime(types.ReactionActionTypeReminder)...)
	reactionActions = append(reactionActions, s.getReactionActionDeleteDailyReminder(types.ReactionActionTypeDailyReminder))
	reactionActions = append(reactionActions, s.getReactionActionDone(types.ReactionActionTypeReminder))
	reactionActions = append(reactionActions, s.getReactionActionReschedule(types.ReactionActionTypeReminder))

	return reactionActions
}

// sendAndStoreMessage is a helper for storing outgoing messages
// it can be used asynchronously
func (s *Syncer) sendAndStoreMessage(message *asyncmessenger.Message, channel *database.Channel, messageType database.MessageType, reminderID uint) {
	response, err := s.messenger.SendMessage(message)
	if err != nil {
		log.Warn("failed sending message with: " + err.Error())
		return
	}

	if messageType != database.MessageTypeDoNotSave {
		reminderIDReal := &reminderID
		if reminderID <= 0 {
			reminderIDReal = nil
		}

		_, err = s.daemon.Database.AddMessage(
			&database.Message{
				Body:               message.Body,
				BodyHTML:           message.BodyHTML,
				Type:               messageType,
				ChannelID:          channel.ID,
				Timestamp:          time.Now().Unix(),
				ExternalIdentifier: response.ExternalIdentifier,
				ReminderID:         reminderIDReal,
			},
		)

		if err != nil {
			log.Warn("Failed to store message to database: " + err.Error())
		}
	}

}

// sendAndStoreReply is a helper for storing outgoing replies
// it can be used asynchronously
func (s *Syncer) sendAndStoreReply(message *asyncmessenger.Response, channel *database.Channel, messageType database.MessageType, reminderID uint) {
	response, err := s.messenger.SendResponse(message)
	if err != nil {
		log.Warn("failed sending message with: " + err.Error())
		return
	}

	if messageType != database.MessageTypeDoNotSave {
		reminderIDReal := &reminderID
		if reminderID <= 0 {
			reminderIDReal = nil
		}

		body, bodyHTML := message.GetResponseMessage()

		_, err = s.daemon.Database.AddMessage(
			&database.Message{
				Body:               body,
				BodyHTML:           bodyHTML,
				Type:               messageType,
				ChannelID:          channel.ID,
				Timestamp:          time.Now().Unix(),
				ExternalIdentifier: response.ExternalIdentifier,
				ReminderID:         reminderIDReal,
				ResponseToMessage:  message.RespondToEventID,
			},
		)

		if err != nil {
			log.Warn("Failed to store message to database: " + err.Error())
		}
	}

}
