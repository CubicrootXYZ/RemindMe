package api_test

import (
	"log/slog"
	"net/http/httptest"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/api"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func testServer(ctrl *gomock.Controller) (*database.MockService, *icaldb.MockService, *httptest.Server) {
	db := database.NewMockService(ctrl)
	icalDB := icaldb.NewMockService(ctrl)

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
