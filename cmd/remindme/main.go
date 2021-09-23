package main

import (
	"log"
	"sync"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/eventdaemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/handler"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/matrixmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/matrixsyncer"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/reminderdaemon"
)

// @title Matrix Reminder and Calendar Bot (RemindMe)
// @version 1.0
// @description API documentation for the matrix reminder and calendar bot. [Inprint & Privacy Policy](https://cubicroot.xyz/impressum) <b> test</b>

// @contact.name Support
// @contact.url https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot

// @host https://your-bot-domain.tld
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.apikey Admin-Authentication
// @in header
// @name Authorization

func main() {
	wg := sync.WaitGroup{}
	// Load config
	config, err := configuration.Load([]string{"config.yml"})
	if err != nil {
		panic(err)
	}

	// Set up database
	db, err := database.Create(config.Database)
	if err != nil {
		panic(err)
	}

	// Create messenger
	messenger, err := matrixmessenger.Create(&config.MatrixBotAccount, db)
	if err != nil {
		panic(err)
	}

	// Create matrix syncer
	syncer := matrixsyncer.Create(config.MatrixBotAccount, config.MatrixUsers, messenger, config.Webserver.BaseURL)

	// Create handler
	calendarHandler := handler.NewCalendarHandler(db)

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
		server := api.NewServer(&config.Webserver, calendarHandler)
		wg.Add(1)
		go server.Start()
	}

	wg.Wait()
	log.Print("Stopped Bot")
}
