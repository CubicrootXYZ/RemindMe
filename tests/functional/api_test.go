package integration

import (
	"net/http"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCalendars(t *testing.T) {
	_, baseURL := tests.NewAPI() // Ensure API is running
	ctx, cancel := tests.ContextWithTimeout()
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/calendar", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "testapikey123456789abcdefg")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	_, status := tests.ParseJSONBodyWithSlice(t, resp.Body)

	assert.Equal(t, "success", status)
	// TODO assert response
}
