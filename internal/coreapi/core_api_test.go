package coreapi_test

import (
	"log/slog"
	"net/http/httptest"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/middleware"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/coreapi"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func testCoreAPI(ctrl *gomock.Controller) (*httptest.Server, *database.MockService) {
	db := database.NewMockService(ctrl)

	api := coreapi.New(&coreapi.Config{
		Database:            db,
		DefaultAuthProvider: middleware.APIKeyAuth("123"),
	}, slog.Default())

	r := gin.New()
	err := api.RegisterRoutes(r)
	if err != nil {
		panic(err)
	}
	server := httptest.NewServer(r)

	return server, db
}

func testDatabaseChannel() database.Channel {
	dailyReminder := uint(130)
	created, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")

	c := database.Channel{
		Description:   "chan desc",
		DailyReminder: &dailyReminder,
	}

	c.ID = 1
	c.CreatedAt = created
	return c
}
