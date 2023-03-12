package coreapi_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/middleware"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/coreapi"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCoreAPI(ctrl *gomock.Controller) (*httptest.Server, *database.MockService) {
	logger := gologger.New(gologger.LogLevelDebug, 0)
	db := database.NewMockService(ctrl)

	api := coreapi.New(&coreapi.Config{
		Database:            db,
		DefaultAuthProvider: middleware.APIKeyAuth("123"),
	}, logger)

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
		TimeZone:      "Europe/Berlin",
	}

	c.ID = 1
	c.CreatedAt = created
	return c
}

func testChannel() coreapi.Channel {
	dailyReminder := "02:10"
	tz := "Europe/Berlin"

	return coreapi.Channel{
		ID:            1,
		CreatedAt:     "2006-01-02T15:04:05+07:00",
		Description:   "chan desc",
		DailyReminder: &dailyReminder,
		TimeZone:      &tz,
	}
}

func TestCoreAPI_ListChannelsHandler(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	server, db := testCoreAPI(ctrl)

	// Mock expectations
	db.EXPECT().GetChannels().Return([]database.Channel{
		testDatabaseChannel(),
	}, nil)

	// Assemble request
	req, err := http.NewRequest(http.MethodGet, server.URL+"/core/channels", nil)
	require.NoError(t, err)

	req.Header.Add("Authorization", "123")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, `{"status":"success","data":[{"ID":1,"CreatedAt":"2006-01-02T15:04:05+07:00","Description":"chan desc","DailyReminder":"02:10","TimeZone":"Europe/Berlin"}]}`, string(body))
}

func TestCoreAPI_ListChannelsHandlerWithError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	server, db := testCoreAPI(ctrl)

	// Mock expectations
	db.EXPECT().GetChannels().Return(nil, errors.New("test"))

	// Assemble request
	req, err := http.NewRequest(http.MethodGet, server.URL+"/core/channels", nil)
	require.NoError(t, err)

	req.Header.Add("Authorization", "123")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, `{"message":"Internal Server Error","status":"error"}`, string(body))
}

func TestCoreAPI_ListChannelsHandlerWithoutAuth(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	server, _ := testCoreAPI(ctrl)

	// Assemble request
	req, err := http.NewRequest(http.MethodGet, server.URL+"/core/channels", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, `{"message":"Unauthenticated","status":"error"}`, string(body))
}
