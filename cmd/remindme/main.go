package main

import (
	"log"
	"sync"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/eventdaemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/matrixmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/matrixsyncer"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/reminderdaemon"
)

func main() {
	wg := sync.WaitGroup{}
	// Load config
	// TODO add timezone settings
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
	messenger, err := matrixmessenger.Create(&config.MatrixBotAccount)
	if err != nil {
		panic(err)
	}

	// Create matrix syncer
	syncer := matrixsyncer.Create(config.MatrixBotAccount, config.MatrixUser)

	// Start event daemon
	eventDaemon := eventdaemon.Create(db, syncer)
	wg.Add(1)
	go eventDaemon.Start(&wg)

	// Start the reminder daemon
	reminderDaemon := reminderdaemon.Create(db, messenger)
	wg.Add(1)
	go reminderDaemon.Start(&wg)

	wg.Wait()
	log.Print("Stopped Bot")
}
