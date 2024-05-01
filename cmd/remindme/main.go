package main

import (
	"os"

	mcmd "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/cmd"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/cmd"
)

// @title Matrix Reminder and Calendar Bot (RemindMe)
// @version 2.0.0
// @description API documentation for the matrix reminder and calendar bot. [Inprint & Privacy Policy](https://cubicroot.xyz/impressum)

// @contact.name Support
// @contact.url https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot

// @host your-bot-domain.tld
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.apikey APIKeyAuthentication
// @in header
// @name Authorization
func main() {
	config, err := cmd.LoadConfiguration()
	if err != nil {
		panic(err)
	}

	config.BuildVersion = mcmd.Version

	err = cmd.Run(config)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
