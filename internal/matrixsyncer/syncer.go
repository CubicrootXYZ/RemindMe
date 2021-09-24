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
	config          configuration.Matrix
	baseURL         string
	adminUsers      []string
	client          *mautrix.Client
	daemon          *eventdaemon.Daemon
	botInfo         *types.BotInfo
	messenger       Messenger
	actions         []*Action         // Actions based on direct messages from the user
	reactionActions []*ReactionAction // Actions based on reactions by the user
	replyActions    []*ReplyAction    // Actions based on replies from the user on existing messages
}

// Create creates a new syncer
func Create(config configuration.Matrix, matrixAdminUsers []string, messenger Messenger, baseURL string) *Syncer {
	syncer := &Syncer{
		config:     config,
		baseURL:    baseURL,
		adminUsers: matrixAdminUsers,
		messenger:  messenger,
	}

	// Add all actions
	syncer.actions = append(syncer.actions, syncer.getActionList())
	syncer.actions = append(syncer.actions, syncer.getActionCommands())
	syncer.actions = append(syncer.actions, syncer.getActionTimezone())
	syncer.actions = append(syncer.actions, syncer.getActionSetDailyReminder())
	syncer.actions = append(syncer.actions, syncer.getActionDeleteDailyReminder())
	syncer.actions = append(syncer.actions, syncer.getActionIcal())
	syncer.actions = append(syncer.actions, syncer.getActionIcalRegenerate())

	syncer.reactionActions = append(syncer.reactionActions, syncer.getReactionActionDelete(ReactionActionTypeReminderRequest))
	syncer.reactionActions = append(syncer.reactionActions, syncer.getReactionsAddTime(ReactionActionTypeReminderRequest)...)
	syncer.reactionActions = append(syncer.reactionActions, syncer.getReactionActionDeleteDailyReminder(ReactionActionTypeDailyReminder))

	syncer.replyActions = append(syncer.replyActions, syncer.getReplyActionDelete(database.MessageTypesWithReminder))
	syncer.replyActions = append(syncer.replyActions, syncer.getReplyActionRecurring(database.MessageTypesWithReminder))

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

	// Initialize handler
	stateMemberHandler := synchandler.NewStateMemberHandler(s.daemon.Database, s.messenger, s.client, s.botInfo)

	// Get messages
	syncer := s.client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, s.handleMessages)
	syncer.OnEventType(event.EventReaction, s.handleReactionEvent)
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
	err := s.daemon.Database.CleanChannels(channels)
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
