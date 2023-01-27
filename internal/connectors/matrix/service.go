package matrix

import (
	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/encryption"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"gorm.io/gorm"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/id"
)

type service struct {
	config   *Config
	logger   gologger.Logger
	database database.Service

	client *mautrix.Client
	crypto struct {
		enabled     bool
		cryptoStore crypto.Store
		stateStore  *encryption.StateStore
	}
}

// Config holds information for the matrix connector.
type Config struct {
	gormDB     *gorm.DB
	Username   string
	Password   string
	Homeserver string
	DeviceID   string
	EnableE2EE bool
	DeviceKey  string
}

// New sets up a new matrix connector.
func New(config *Config, database *database.Service, logger gologger.Logger) (Service, error) {
	logger.Debugf("setting up matrix connector ...")
	service := &service{
		config:   config,
		logger:   logger,
		database: *database,
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

	logger.Debugf("matrix connector setup finished")
	return service, nil
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

	stateStore := encryption.NewStateStore(nil, &encryption.StateStoreConfig{ // TODO database
		Username:   service.config.Username,
		Homeserver: service.config.Homeserver,
	}, service.logger.WithField("component", "statestore"))
	service.crypto.stateStore = stateStore

	service.config.DeviceID = deviceID.String() // we might get a new device ID if none is set
	service.crypto.enabled = true

	service.logger.Debugf("matrix end to end encryption setup finished")
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
