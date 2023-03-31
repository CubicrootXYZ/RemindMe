package coreapi_test

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

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
	assert.Equal(t, `{"status":"success","data":[{"ID":1,"CreatedAt":"2006-01-02T15:04:05+07:00","Description":"chan desc","DailyReminder":"02:10"}]}`, string(body))
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
