package tests

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/handler"
)

var testServer *api.Server

// NewAPI builds a new API for testing.
// Makes a few assumptions about a test database and ports. Check out the code for details.
// Singleton, ensure your tests are able to run in parallel.
func NewAPI() (*api.Server, string) {
	if testServer != nil {
		return testServer, "http://localhost:4232"
	}

	db := NewDatabase()

	testServer = api.NewServer(
		&configuration.Webserver{
			Enabled: true,
			APIkey:  "testapikey123456789abcdefg",
			BaseURL: "localhost:4232",
			Address: ":4232",
		},
		handler.NewCalendarHandler(db),
		handler.NewDatabaseHandler(db),
	)

	go testServer.Start(true)

	time.Sleep(time.Second) // Wait for server to start

	return testServer, "http://localhost:4232"
}
