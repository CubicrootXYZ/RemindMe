package matrix

import (
	"errors"
	"regexp"
	"time"

	"github.com/CubicrootXYZ/gologger"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/encryption"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/id"
)

type service struct {
	config         *Config
	logger         gologger.Logger
	database       database.Service
	matrixDatabase matrixdb.Service
	messenger      messenger.Messenger
	botname        string

	client *mautrix.Client
	crypto struct {
		enabled     bool
		cryptoStore crypto.Store
		stateStore  *encryption.StateStore
		olm         *crypto.OlmMachine
	}
	lastMessageFrom time.Time
}

//go:generate mockgen -destination=message_action_mock.go -package=matrix . MessageAction

// MessageAction defines an interface for an action on messages.
type MessageAction interface {
	Selector() *regexp.Regexp
	Name() string
	HandleEvent(event *MessageEvent)
	Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, bridgeServices *BridgeServices)
}

//go:generate mockgen -destination=reply_action_mock.go -package=matrix . ReplyAction

// ReplyAction defines an interface for an action on replies.
type ReplyAction interface {
	Selector() *regexp.Regexp
	Name() string
	HandleEvent(event *MessageEvent, replyToMessage *matrixdb.MatrixMessage)
	Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, bridgeServices *BridgeServices)
}

// Config holds information for the matrix connector.
type Config struct {
	Username   string
	Password   string
	Homeserver string
	DeviceID   string
	EnableE2EE bool
	DeviceKey  string

	MessageActions       []MessageAction
	DefaultMessageAction MessageAction
	ReplyActions         []ReplyAction
	DefaultReplyAction   ReplyAction

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
func New(config *Config, database database.Service, matrixDB matrixdb.Service, logger gologger.Logger) (Service, error) {
	logger.Debugf("setting up matrix connector ...")

	service := &service{
		config:         config,
		logger:         logger,
		database:       database,
		matrixDatabase: matrixDB,
		botname:        format.FullUsername(config.Username, config.Homeserver),
	}

	err := service.setupMautrixClient()
	if err != nil {
		service.logger.Err(err)
		return nil, err
	}

	if config.EnableE2EE {
		err = service.setupEncryption()
		if err != nil {
			service.logger.Err(err)
			return nil, err
		}
	}

	err = service.setupMessenger() // important to call after crypto setup
	if err != nil {
		return nil, err
	}

	service.setupActions() // important to call after crypto, db and mautrix setup

	service.setLastMessage()

	logger.Debugf("matrix connector setup finished")
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
	actions := []interface {
		Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, bridgeServices *BridgeServices)
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

	for _, action := range actions {
		action.Configure(
			service.logger.WithField("component", "action-"+action.Name()),
			service.client,
			service.messenger,
			service.matrixDatabase,
			service.database,
			service.config.BridgeServices,
		)
	}
}

func (service *service) setupMautrixClient() error {
	service.logger.Debugf("setting up mautrix client ...")

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

	service.logger.Debugf("matrix client setup finished")
	return err
}

func (service *service) setupEncryption() error {
	service.logger.Debugf("setting up matrix end to end encryption ...")

	cryptoStore, deviceID, err := encryption.NewCryptoStore(
		service.config.Username,
		service.config.DeviceKey,
		service.config.Homeserver,
		service.config.DeviceID,
		service.logger.WithField("component", "cryptostore"),
	)
	if err != nil {
		return err
	}
	service.crypto.cryptoStore = cryptoStore

	stateStore := encryption.NewStateStore(service.matrixDatabase, &encryption.StateStoreConfig{
		Username:   service.config.Username,
		Homeserver: service.config.Homeserver,
	}, service.logger.WithField("component", "statestore"))
	service.crypto.stateStore = stateStore

	olm, err := encryption.NewOlmMachine(service.client, service.crypto.cryptoStore, service.crypto.stateStore, service.logger.WithField("component", "olm"))
	if err != nil {
		service.logger.Errorf("failed setting up olm machine: %s", err.Error())
		return err
	}
	service.crypto.olm = olm

	service.config.DeviceID = deviceID.String() // we might get a new device ID if none is set
	service.crypto.enabled = true

	service.logger.Debugf("matrix end to end encryption setup finished")
	return nil
}

func (service *service) setupMessenger() error {
	config := &messenger.Config{}

	if service.crypto.enabled {
		config.Crypto = &messenger.CryptoTools{}
		config.Crypto.Olm = service.crypto.olm
		config.Crypto.StateStore = service.crypto.stateStore
	}

	messenger, err := messenger.NewMessenger(config, service.matrixDatabase, service.client,
		service.logger.WithField("component", "matrix-messenger"))
	if err != nil {
		return err
	}

	service.messenger = messenger

	return nil
}

// Start starts the services asynchronous processes.
// This method will block until stopped.
func (service *service) Start() error {
	service.logger.Debugf("starting matrix connector")
	err := service.startListener()
	service.logger.Debugf("matrix connector stopped")
	return err
}

// Stop stops the services asynchronous processes.
// This method will not block, wait for Stop() to return.
func (service *service) Stop() error {
	service.logger.Debugf("stopping matrix connector ...")
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
	if err != nil && !errors.Is(err, matrixdb.ErrNotFound) {
		return err
	}

	err = service.removeRoom(room)
	return err
}

// OutputRemoved to tell the connector an output got removed.
func (service *service) OutputRemoved(outputType string, outputID uint) error {
	if outputType != OutputType {
		// Input is not from this connector, ignore
		return nil
	}

	room, err := service.matrixDatabase.GetRoomByID(outputID)
	if err != nil && !errors.Is(err, matrixdb.ErrNotFound) {
		return err
	}

	err = service.removeRoom(room)
	return err
}
