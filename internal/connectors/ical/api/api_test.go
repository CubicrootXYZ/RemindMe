package api_test

import (
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/api"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
)

func testServer(t *testing.T) (*database.MockService, *icaldb.MockService, *httptest.Server) {
	t.Helper()
	db := database.NewMockService(t)
	icalDB := icaldb.NewMockService(t)

	api := api.New(&api.Config{
		Database: db,
		IcalDB:   icalDB,
	}, slog.Default())

	r := gin.New()

	err := api.RegisterRoutes(r)
	if err != nil {
		panic(err)
	}

	return db, icalDB, httptest.NewServer(r)
}
