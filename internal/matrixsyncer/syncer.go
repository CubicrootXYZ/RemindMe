package matrixsyncer

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/eventdaemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
)

// Syncer receives messages from a matrix channel
type Syncer struct {
	config    configuration.Matrix
	user      string
	client    *mautrix.Client
	daemon    *eventdaemon.Daemon
	botName   string
	messenger Messenger
}

// Create creates a new syncer
func Create(config configuration.Matrix, matrixUser string, messenger Messenger) *Syncer {
	return &Syncer{
		config:    config,
		user:      matrixUser,
		messenger: messenger,
	}
}

// Start starts the syncer
func (s *Syncer) Start(daemon *eventdaemon.Daemon) error {
	log.Info(fmt.Sprintf("Starting Matrix syncer for user %s on server %s", s.config.Username, s.config.Homeserver))

	s.daemon = daemon
	s.botName = fmt.Sprintf("@%s:%s", s.config.Username, strings.ReplaceAll(strings.ReplaceAll(s.config.Homeserver, "https://", ""), "http://", ""))

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

	_, err = s.daemon.Database.GetChannelByUserIdentifier(s.user)
	if err == gorm.ErrRecordNotFound {
		_, err = s.createChannel(s.user)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Get messages
	syncer := s.client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, s.handleMessages)
	return client.Sync()
}

// Stop stops the syncer
func (s *Syncer) Stop() {
	s.client.StopSync()
}
