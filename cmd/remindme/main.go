package main

import (
	"os"
	"sync"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/encryption"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/eventdaemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/handler"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/matrixmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/matrixsyncer"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/reminderdaemon"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/id"
)

// @title Matrix Reminder and Calendar Bot (RemindMe)
// @version 1.5.2
// @description API documentation for the matrix reminder and calendar bot. [Inprint & Privacy Policy](https://cubicroot.xyz/impressum)

// @contact.name Support
// @contact.url https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot

// @host your-bot-domain.tld
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.apikey Admin-Authentication
// @in header
// @name Authorization
func main() {
	for {
		err := startup()

		if err == nil {
			log.Info("Bot stopped cleanly - exiting")
			break
		}

		log.Info("Bot stopped due to error: " + err.Error())
		log.Info("Will retry in 3 minutes")
		time.Sleep(3 * time.Minute)
	}
}

func startup() error {
	log.Info("Starting up bot")
	wg := sync.WaitGroup{}

	// Make data directory
	err := os.MkdirAll("data", 0755)
	if err != nil {
		return err
	}

	// Load config
	config, err := configuration.Load([]string{"config.yml"})
	if err != nil {
		return err
	}

	// Initialize logger
	logger := log.InitLogger(config.Debug)
	defer logger.Sync()

	// Set up database
	log.Debug("Initializing database")
	db, err := database.Create(config.Database, config.Debug)
	if err != nil {
		return err
	}

	// Create encryption handler
	cryptoStore, stateStore, _, err := initializeEncryptionSetup(&config.MatrixBotAccount, db, config.Debug)
	if err != nil {
		return err
	}

	// Create matrix client
	matrixClient, err := initializeMatrixClient(&config.MatrixBotAccount)
	if err != nil {
		return err
	}

	// Inject matrix client into database
	db.SetMatrixClient(matrixClient)

	// Create messenger
	log.Debug("Creating messenger")
	messenger, err := matrixmessenger.Create(config.Debug, config, db, cryptoStore, stateStore, matrixClient)
	if err != nil {
		return err
	}

	// Create matrix syncer
	log.Debug("Creating syncer and handlers")
	syncer := matrixsyncer.Create(config, config.MatrixUsers, messenger, cryptoStore, stateStore, matrixClient)

	// Create handler
	calendarHandler := handler.NewCalendarHandler(db)
	databaseHandler := handler.NewDatabaseHandler(db)

	// Start event daemon
	log.Debug("Starting up event daemon")
	eventDaemon := eventdaemon.Create(db, syncer)
	wg.Add(1)
	go eventDaemon.Start(&wg)

	// Start the reminder daemon
	log.Debug("Starting up reminder daemon")
	reminderDaemon := reminderdaemon.Create(db, messenger)
	wg.Add(1)
	go reminderDaemon.Start(&wg)

	// Start the Webserver
	if config.Webserver.Enabled {
		log.Debug("Starting up webserver")
		server := api.NewServer(&config.Webserver, calendarHandler, databaseHandler)
		wg.Add(1)
		go server.Start(config.Debug)
	}

	log.Info("Started successfully :)")
	wg.Wait()
	log.Info("Stopped Bot")

	return nil
}

func initializeEncryptionSetup(config *configuration.Matrix, db *database.Database, debug bool) (cryptoStore crypto.Store, stateStore *encryption.StateStore, deviceID id.DeviceID, err error) {
	log.Debug("Initializing encryption setup ...")

	deviceID = id.DeviceID(config.DeviceID)

	sqlDB, err := db.SQLDB()
	if err != nil {
		return
	}
	if config.E2EE {
		cryptoStore, deviceID, err = encryption.GetCryptoStore(debug, sqlDB, config)
		if err != nil {
			return
		}
		stateStore = encryption.NewStateStore(db, config)
		config.DeviceID = deviceID.String()
	}

	log.Debug("... finished initializing encryption setup")

	return
}

func initializeMatrixClient(config *configuration.Matrix) (matrixClient *mautrix.Client, err error) {
	log.Debug("Initializing matrix client ...")

	matrixClient, err = mautrix.NewClient(config.Homeserver, "", "")
	if err != nil {
		return
	}

	_, err = matrixClient.Login(&mautrix.ReqLogin{
		Type:             "m.login.password",
		Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: config.Username},
		Password:         config.Password,
		DeviceID:         id.DeviceID(config.DeviceID),
		StoreCredentials: true,
	})
	if err != nil {
		return
	}

	log.Debug("... finished initializing matrix client")
	return
}
