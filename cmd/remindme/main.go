package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/encryption"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/eventdaemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/handler"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/icalimporter"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/matrixsyncer"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/reminderdaemon"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/id"
)

// @title Matrix Reminder and Calendar Bot (RemindMe)
// @version 1.8.0
// @description API documentation for the matrix reminder and calendar bot. [Inprint & Privacy Policy](https://cubicroot.xyz/impressum)

// @contact.name Support
// @contact.url https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot

// @host your-bot-domain.tld
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.apikey AdminAuthentication
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

	config, logger, db, err := setupBasics()
	if err != nil {
		return err
	}

	eventDaemon, reminderDaemon, handler, err := initializeDependencies(config, db)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(context.Background())

	// Start event daemon
	log.Debug("starting up event daemon")
	eg.Go(func() error {
		return eventDaemon.Start()
	})

	// Start the reminder daemon
	log.Debug("starting up reminder daemon")
	eg.Go(func() error {
		return reminderDaemon.Start()
	})

	// Start the Webserver
	var server *api.Server
	if config.Webserver.Enabled {
		log.Debug("starting up webserver")
		server = api.NewServer(&config.Webserver, handler)
		eg.Go(func() error {
			return server.Start(config.Debug)
		})
	}

	// Start the ical importer
	var icalImporter icalimporter.IcalImporter
	if config.BotSettings.AllowIcalImport {
		log.Debug("starting up ical importer")
		icalImporter = icalimporter.NewIcalImporter(db)
		eg.Go(func() error {
			return icalImporter.Run()
		})
	}

	// Listen to signals and shut down if receiving signal
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	shutdown := make(chan interface{})

	go func() {
		logger.Info("waiting for signal")
		select {
		case s := <-sigc:
			logger.Info("got signal, shutting down: ", s)
		case <-ctx.Done():
			logger.Info("a process stopped, shutting down")
		}
		shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		// Shut down all routines
		if icalImporter != nil {
			icalImporter.Stop()
		}
		if server != nil {
			err := server.Stop(shutdownCtx)
			if err != nil && err != http.ErrServerClosed {
				logger.Error(err.Error())
			}
		}
		reminderDaemon.Stop()
		eventDaemon.Stop()

		close(shutdown)
	}()

	log.Info("started successfully :)")
	err = eg.Wait()
	if err != nil {
		logger.Error("process failed with: ", err.Error())
	}
	log.Info("stopped bot")

	<-shutdown

	return nil
}

func setupBasics() (*configuration.Config, *zap.SugaredLogger, *database.Database, error) {
	// Make data directory
	err := os.MkdirAll("data", 0755)
	if err != nil {
		return nil, nil, nil, err
	}

	// Load config
	config, err := configuration.Load([]string{"config.yml"})
	if err != nil {
		return nil, nil, nil, err
	}

	// Initialize logger
	logger := log.InitLogger(config.Debug)

	// Set up database
	log.Debug("Initializing database")
	db, err := database.Create(config.Database, config.Debug)
	if err != nil {
		return nil, nil, nil, err
	}

	return config, logger, db, nil
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

func initializeDependencies(config *configuration.Config, db *database.Database) (*eventdaemon.Daemon, *reminderdaemon.Daemon, *api.Handler, error) {
	// Create encryption handler
	cryptoStore, stateStore, _, err := initializeEncryptionSetup(&config.MatrixBotAccount, db, config.Debug)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create matrix client
	matrixClient, err := initializeMatrixClient(&config.MatrixBotAccount)
	if err != nil {
		return nil, nil, nil, err
	}

	// Inject matrix client into database
	db.SetMatrixClient(matrixClient)

	// Create messenger
	log.Debug("creating messenger")
	messenger, err := asyncmessenger.NewMessenger(config.Debug, config, db, cryptoStore, stateStore, matrixClient)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create matrix syncer
	log.Debug("creating syncer and handlers")
	syncer := matrixsyncer.Create(config, config.MatrixUsers, messenger, cryptoStore, stateStore, matrixClient)

	// Create handler
	calendarHandler := handler.NewCalendarHandler(db)
	databaseHandler := handler.NewDatabaseHandler(db)

	eventDaemon := eventdaemon.Create(db, syncer)
	reminderDaemon := reminderdaemon.Create(db, messenger)

	return eventDaemon,
		reminderDaemon,
		&api.Handler{
			Calendar: calendarHandler,
			Database: databaseHandler,
		},
		nil
}
