package matrix

import (
	"errors"
	"log/slog"
	"regexp"
	"time"

	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

var (
	metricEventInCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "remindme",
		Name:      "matrix_events_in_total",
		Help:      "Counts events received by the matrix connector.",
	}, []string{"event_type"})
)

type service struct {
	config         *Config
	logger         *slog.Logger
	database       database.Service
	matrixDatabase matrixdb.Service
	messenger      messenger.Messenger
	botname        string

	client          *mautrix.Client
	lastMessageFrom time.Time

	metricEventInCount *prometheus.CounterVec
}

//go:generate mockgen -destination=message_action_mock.go -package=matrix . MessageAction

// MessageAction defines an interface for an action on messages.
type MessageAction interface {
	Selector() *regexp.Regexp
	Name() string
	HandleEvent(event *MessageEvent)
	Configure(logger *slog.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, bridgeServices *BridgeServices)
}

//go:generate mockgen -destination=reply_action_mock.go -package=matrix . ReplyAction

// ReplyAction defines an interface for an action on replies.
type ReplyAction interface {
	Selector() *regexp.Regexp
	Name() string
	HandleEvent(event *MessageEvent, replyToMessage *matrixdb.MatrixMessage)
	Configure(logger *slog.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, bridgeServices *BridgeServices)
}

//go:generate mockgen -destination=reaction_action_mock.go -package=matrix . ReactionAction

// ReactionAction defines an Interface for an action on reactions.
type ReactionAction interface {
	Selector() []string
	Name() string
	HandleEvent(event *ReactionEvent, reactionToMessage *matrixdb.MatrixMessage)
	Configure(logger *slog.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, bridgeServices *BridgeServices)
}

// Config holds information for the matrix connector.
type Config struct {
	Username   string
	Password   string
	Homeserver string
	DeviceID   string
	DeviceKey  string

	MessageActions       []MessageAction
	DefaultMessageAction MessageAction
	ReplyActions         []ReplyAction
	DefaultReplyAction   ReplyAction
	ReactionActions      []ReactionAction

	AllowInvites  bool
	RoomLimit     uint
	UserWhitelist []string // Invites frim this users will allways be followed

	BridgeServices *BridgeServices
}

// BridgeServices contains services where the matrix connector acts as a bridge, e.g.
// because they do not have any user interface.
type BridgeServices struct {
	ICal BridgeServiceICal
}

// BridgeServiceICal is an interface for a bridge to the iCal connector.
type BridgeServiceICal interface {
	NewOutput(channelID uint) (*icaldb.IcalOutput, string, error)
	GetOutput(outputID uint, regenToken bool) (*icaldb.IcalOutput, string, error)
}

// New sets up a new matrix connector.
func New(config *Config, database database.Service, matrixDB matrixdb.Service, logger *slog.Logger) (Service, error) {
	logger.Debug("setting up matrix connector")

	service := &service{
		config:         config,
		logger:         logger,
		database:       database,
		matrixDatabase: matrixDB,
		botname:        format.FullUsername(config.Username, config.Homeserver),

		metricEventInCount: metricEventInCount,
	}

	err := service.setupMautrixClient()
	if err != nil {
		service.logger.Error("failed to setup mautrix client", "error", err)
		return nil, err
	}

	err = service.setupMessenger()
	if err != nil {
		return nil, err
	}

	service.setupActions() // important to call after db and mautrix setup

	service.setLastMessage()

	logger.Debug("matrix connector setup finished")
	return service, nil
}

// setLastMessage so the handlers will know which messages can be ignored savely
func (service *service) setLastMessage() {
	message, err := service.matrixDatabase.GetLastMessage()
	if err == nil {
		service.lastMessageFrom = message.SendAt
	}

	event, err := service.matrixDatabase.GetLastEvent()
	if err == nil && event.SendAt.Sub(service.lastMessageFrom) > 0 {
		service.lastMessageFrom = event.SendAt
	}
}

func (service *service) setupActions() {
	// Collect all actions and inject dependencies.
	actions := []interface {
		Configure(logger *slog.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, bridgeServices *BridgeServices)
		Name() string
	}{
		service.config.DefaultMessageAction,
		service.config.DefaultReplyAction,
	}

	for _, action := range service.config.MessageActions {
		actions = append(actions, action)
	}
	for _, action := range service.config.ReplyActions {
		actions = append(actions, action)
	}
	for _, action := range service.config.ReactionActions {
		actions = append(actions, action)
	}

	for _, action := range actions {
		action.Configure(
			service.logger.With("component", "action-"+action.Name()),
			service.client,
			service.messenger,
			service.matrixDatabase,
			service.database,
			service.config.BridgeServices,
		)
	}
}

func (service *service) setupMautrixClient() error {
	service.logger.Debug("setting up mautrix client")

	matrixClient, err := mautrix.NewClient(service.config.Homeserver, "", "")
	if err != nil {
		return err
	}

	service.client = matrixClient

	_, err = service.client.Login(&mautrix.ReqLogin{
		Type:             "m.login.password",
		Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: service.config.Username},
		Password:         service.config.Password,
		DeviceID:         id.DeviceID(service.config.DeviceID),
		StoreCredentials: true,
	})

	service.logger.Debug("mautrix client setup finished")
	return err
}

func (service *service) setupMessenger() error {
	config := &messenger.Config{}

	messenger, err := messenger.NewMessenger(config, service.matrixDatabase, service.client,
		service.logger.With("component", "matrix-messenger"))
	if err != nil {
		return err
	}

	service.messenger = messenger

	return nil
}

// Start starts the services asynchronous processes.
// This method will block until stopped.
func (service *service) Start() error {
	service.logger.Debug("starting matrix connector")
	err := service.startListener()
	service.logger.Debug("matrix connector stopped")
	return err
}

// Stop stops the services asynchronous processes.
// This method will not block, wait for Stop() to return.
func (service *service) Stop() error {
	service.logger.Debug("stopping matrix connector ...")
	service.client.StopSync()
	return nil
}

// InputRemoved to tell the connector an input got removed.
func (service *service) InputRemoved(inputType string, inputID uint) error {
	if inputType != InputType {
		// Input is not from this connector, ignore
		return nil
	}

	room, err := service.matrixDatabase.GetRoomByID(inputID)
	if err != nil {
		if errors.Is(err, matrixdb.ErrNotFound) {
			return nil
		}
		return err
	}

	err = service.removeRoom(room)
	return err
}

// OutputRemoved to tell the connector an output got removed.
func (service *service) OutputRemoved(outputType string, outputID uint) error {
	if outputType != OutputType {
		// Output is not from this connector, ignore
		return nil
	}

	room, err := service.matrixDatabase.GetRoomByID(outputID)
	if err != nil {
		if errors.Is(err, matrixdb.ErrNotFound) {
			// No need to delete already deleted room.
			return nil
		}
		return err
	}

	err = service.removeRoom(room)
	return err
}
