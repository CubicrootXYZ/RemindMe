package main

import (
	pLog "log"
	"os"
	"sync"

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
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/id"
)

// @title Matrix Reminder and Calendar Bot (RemindMe)
// @version 1.3.0
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
	wg := sync.WaitGroup{}

	// Make data directory
	err := os.MkdirAll("data", 0755)
	if err != nil {
		panic(err)
	}

	// Load config
	config, err := configuration.Load([]string{"config.yml"})
	if err != nil {
		panic(err)
	}

	logger := log.InitLogger(config.Debug)
	defer logger.Sync()

	// Set up database
	db, err := database.Create(config.Database, config.Debug)
	if err != nil {
		panic(err)
	}

	// Create encryption handler
	var cryptoStore crypto.Store
	var stateStore *encryption.StateStore
	deviceID := id.DeviceID(config.MatrixBotAccount.DeviceID) //lint:ignore SA4006 Needed as backup here

	sqlDB, err := db.SQLDB()
	if err != nil {
		panic(err)
	}
	if config.MatrixBotAccount.E2EE {
		cryptoStore, deviceID, err = encryption.GetCryptoStore(sqlDB, &config.MatrixBotAccount)
		if err != nil {
			panic(err)
		}
		stateStore = encryption.NewStateStore(db)
		config.MatrixBotAccount.DeviceID = deviceID.String()
	}

	// Create messenger
	messenger, err := matrixmessenger.Create(&config.MatrixBotAccount, db, cryptoStore, stateStore)
	if err != nil {
		panic(err)
	}

	// Create matrix syncer
	syncer := matrixsyncer.Create(config, config.MatrixUsers, messenger, cryptoStore, stateStore)

	// Create handler
	calendarHandler := handler.NewCalendarHandler(db)
	databaseHandler := handler.NewDatabaseHandler(db)

	// Start event daemon
	eventDaemon := eventdaemon.Create(db, syncer)
	wg.Add(1)
	go eventDaemon.Start(&wg)

	// Start the reminder daemon
	reminderDaemon := reminderdaemon.Create(db, messenger)
	wg.Add(1)
	go reminderDaemon.Start(&wg)

	// Start the Webserver
	if config.Webserver.Enabled {
		server := api.NewServer(&config.Webserver, calendarHandler, databaseHandler)
		wg.Add(1)
		go server.Start(config.Debug)
	}

	wg.Wait()
	pLog.Print("Stopped Bot")
}
