package matrixsyncer

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/eventdaemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/synchandler"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
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
	messenger   Messenger
}

// Create creates a new syncer
func Create(config *configuration.Config, matrixAdminUsers []string, messenger Messenger) *Syncer {
	syncer := &Syncer{
		config:      config.MatrixBotAccount,
		baseURL:     config.Webserver.BaseURL,
		adminUsers:  matrixAdminUsers,
		messenger:   messenger,
		botSettings: &config.BotSettings,
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

	// Log into matrix
	client, err := mautrix.NewClient(s.config.Homeserver, "", "")
	if err != nil {
		return err
	}

	s.client = client
	_, err = s.client.Login(&mautrix.ReqLogin{
		Type:             "m.login.password",
		Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: s.config.Username},
		Password:         s.config.Password,
		StoreCredentials: true,
	})
	if err != nil {
		return err
	}
	log.Info("Logged in to matrix")

	err = s.syncChannels()
	if err != nil {
		return err
	}

	// Load actions
	messageActions := s.getActions()
	replyActions := s.getReplyActions()
	reactionActions := s.getReactionActions()

	// Initialize handler
	messageHandler := synchandler.NewMessageHandler(s.daemon.Database, s.messenger, s.botInfo, replyActions, messageActions)
	stateMemberHandler := synchandler.NewStateMemberHandler(s.daemon.Database, s.messenger, s.client, s.botInfo, s.botSettings)
	reactionHandler := synchandler.NewReactionHandler(s.daemon.Database, s.messenger, s.botInfo, reactionActions)

	// Get messages
	syncer := s.client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, messageHandler.NewEvent)
	syncer.OnEventType(event.EventReaction, reactionHandler.NewEvent)
	syncer.OnEventType(event.StateMember, stateMemberHandler.NewEvent)
	return client.Sync()
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

		s.messenger.SendNotice("Sorry I was sleeping for a while. I am now ready for your requests!", channel.ChannelIdentifier)
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
			panic(err)
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
	reactionActions = append(reactionActions, s.getReactionsAddTime(types.ReactionActionTypeReminderRequest)...)
	reactionActions = append(reactionActions, s.getReactionActionDeleteDailyReminder(types.ReactionActionTypeDailyReminder))

	return reactionActions
}
