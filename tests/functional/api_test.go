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

	data, status := tests.ParseJSONBodyWithSlice(t, resp.Body)

	assert.Equal(t, "success", status)

	for _, channelRaw := range data {
		channel, ok := channelRaw.(map[string]interface{})
		require.True(t, ok, "type casting channel failed")

		switch channel["id"] {
		case float64(1):
			assert.Equal(t, "!123456789", channel["channel_id"])
			assert.Equal(t, "testuser@example.com", channel["user_id"])
			assert.NotEmpty(t, channel["token"])
		case float64(2):
			assert.Equal(t, "!abcdefghij", channel["channel_id"])
			assert.Equal(t, "testuser2@example.com", channel["user_id"])
			assert.NotEmpty(t, channel["token"])
		default:
			require.Failf(t, "unknown channel", "unknown channel id: %v (%T)", channel["id"], channel["id"])
		}
	}
}
