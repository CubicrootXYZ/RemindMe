package tests

import (
	"context"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/handler"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

var testServer *api.Server
var router routers.Router

// NewAPI builds a new API for testing.
// Makes a few assumptions about a test database and ports. Check out the code for details.
// Singleton, ensure your tests are able to run in parallel.
func NewAPI() (*api.Server, string, routers.Router) {
	if testServer != nil {
		return testServer, "http://localhost:4232", router
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

	router = loadOpenAPISpecs()

	time.Sleep(time.Second) // Wait for server to start

	return testServer, "http://localhost:4232", router
}

func loadOpenAPISpecs() routers.Router {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile("../../docs/api-spec.yaml")
	if err != nil {
		panic(err)
	}

	doc.Servers = nil // ignores domain in validation

	// Validate document
	err = doc.Validate(ctx)
	if err != nil {
		panic(err)
	}

	specRouter, err := gorillamux.NewRouter(doc)
	if err != nil {
		panic(err)
	}

	return specRouter
}
