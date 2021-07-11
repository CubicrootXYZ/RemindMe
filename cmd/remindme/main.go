package main

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/eventdaemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/matrixsyncer"
)

func main() {
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

	// Create matrix syncer
	syncer := matrixsyncer.Create(config.MatrixBotAccount, config.MatrixUser)

	// Start event daemon
	eventDaemon := eventdaemon.Create(db, syncer)
	eventDaemon.Start()

	// TODO add a reminder daemon
}
