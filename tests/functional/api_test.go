package integration

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCalendar(t *testing.T) {
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

	for _, calendarRaw := range data {
		calendar, ok := calendarRaw.(map[string]interface{})
		require.True(t, ok, "type casting channel failed")

		switch calendar["id"] {
		case float64(1):
			assert.Equal(t, "!123456789", calendar["channel_id"])
			assert.Equal(t, "testuser@example.com", calendar["user_id"])
			assert.NotEmpty(t, calendar["token"])
		case float64(2):
			assert.Equal(t, "!abcdefghij", calendar["channel_id"])
			assert.Equal(t, "testuser2@example.com", calendar["user_id"])
			assert.NotEmpty(t, calendar["token"])
		case float64(3):
			// Channels added later but are not required for this test.
		default:
			require.Failf(t, "unknown calendar", "unknown clanedar with id: %v (%T)", calendar["id"], calendar["id"])
		}
	}
}

func TestPatchCalendarID(t *testing.T) {
	_, baseURL := tests.NewAPI() // Ensure API is running
	ctx, cancel := tests.ContextWithTimeout()
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, baseURL+"/calendar/1", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "testapikey123456789abcdefg")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	_, status := tests.ParseJSONBodyWithMessage(t, resp.Body)
	assert.Equal(t, "success", status)
}

func TestGetCalendarIDIcal(t *testing.T) {
	_, baseURL := tests.NewAPI() // Ensure API is running
	ctx, cancel := tests.ContextWithTimeout()
	defer cancel()

	// Get token
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/calendar", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "testapikey123456789abcdefg")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, status := tests.ParseJSONBodyWithSlice(t, resp.Body)
	require.Equal(t, "success", status)

	token := data[0].(map[string]interface{})["token"].(string)
	id := data[0].(map[string]interface{})["id"].(float64)
	require.NotEmpty(t, token)
	require.Greater(t, id, 0.0)

	// Get ICAL
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/calendar/%d/ical?token=%s", baseURL, int(id), token), nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:RemindMe\nMETHOD:PUBLISH\nEND:VCALENDAR\n", string(body))
}

func TestGetChannel(t *testing.T) {
	_, baseURL := tests.NewAPI() // Ensure API is running
	ctx, cancel := tests.ContextWithTimeout()
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/channel", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "testapikey123456789abcdefg")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, status := tests.ParseJSONBodyWithSlice(t, resp.Body)
	assert.Equal(t, "success", status)

	for _, calendarRaw := range data {
		calendar, ok := calendarRaw.(map[string]interface{})
		require.True(t, ok, "type casting channel failed")

		switch calendar["id"] {
		case float64(1):
			assert.Equal(t, "!123456789", calendar["channel_id"])
			assert.Equal(t, "testuser@example.com", calendar["user_id"])
			assert.Empty(t, calendar["timezone"])
			assert.True(t, calendar["daily_reminder"].(bool))
		case float64(2):
			assert.Equal(t, "!abcdefghij", calendar["channel_id"])
			assert.Equal(t, "testuser2@example.com", calendar["user_id"])
			assert.Equal(t, "admin", calendar["role"])
			assert.Equal(t, "Berlin", calendar["timezone"])
			assert.True(t, calendar["daily_reminder"].(bool))
		case float64(3):
			// Channels that were added later, but are not needed for this test.
		default:
			require.Failf(t, "unknown channel", "unknown channel id: %v (%T)", calendar["id"], calendar["id"])
		}

		assert.NotEmpty(t, calendar["created"])
	}
}

func TestDeleteChannelID(t *testing.T) {
	_, baseURL := tests.NewAPI() // Ensure API is running
	ctx, cancel := tests.ContextWithTimeout()
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, baseURL+"/channel/3", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "testapikey123456789abcdefg")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	_, status := tests.ParseJSONBodyWithMessage(t, resp.Body)
	assert.Equal(t, "success", status)
}

func TestChannelIDThirdPartyResources(t *testing.T) {
	// Test Add, List and Delete in one go
	_, baseURL := tests.NewAPI() // Ensure API is running
	ctx, cancel := tests.ContextWithTimeout()
	defer cancel()

	// Add
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/channel/1/thirdpartyresources", strings.NewReader(`{"type":"ical", "url":"testurl"}`))
	require.NoError(t, err)
	req.Header.Add("Authorization", "testapikey123456789abcdefg")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	_, status := tests.ParseJSONBodyWithMessage(t, resp.Body)
	assert.Equal(t, "success", status)

	// List
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/channel/1/thirdpartyresources", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "testapikey123456789abcdefg")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, status := tests.ParseJSONBodyWithSlice(t, resp.Body)
	assert.Equal(t, "success", status)

	for _, resourceRaw := range data {
		resource := resourceRaw.(map[string]interface{})

		switch resource["id"] {
		case float64(1):
			assert.Equal(t, "testurl", resource["url"])
			assert.Equal(t, "ICAL", resource["type"])
		default:
			require.Failf(t, "unknown resource", "unknown resource id: %v (%T)", resource["id"], resource["id"])
		}
	}

	// Delete
	req, err = http.NewRequestWithContext(ctx, http.MethodDelete, baseURL+"/channel/1/thirdpartyresources/1", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "testapikey123456789abcdefg")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	_, status = tests.ParseJSONBodyWithMessage(t, resp.Body)
	assert.Equal(t, "success", status)
}

func TestGetUser(t *testing.T) {
	_, baseURL := tests.NewAPI() // Ensure API is running
	ctx, cancel := tests.ContextWithTimeout()
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/user", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "testapikey123456789abcdefg")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, status := tests.ParseJSONBodyWithSlice(t, resp.Body)
	assert.Equal(t, "success", status)

	for _, userRaw := range data {
		user := userRaw.(map[string]interface{})

		switch user["user_id"] {
		case "testuser@example.com", "testuser2@example.com", "testuser3@example.com":
			assert.Equal(t, false, user["blocked"].(bool))
			assert.Empty(t, user["comment"])
			assert.Equal(t, 1, len(user["channels"].([]interface{})))
		default:
			require.Failf(t, "unknown user", "unknown user id: %v (%T)", user["user_id"], user["user_id"])
		}
	}
}

func TestPutUserID(t *testing.T) {
	_, baseURL := tests.NewAPI() // Ensure API is running
	ctx, cancel := tests.ContextWithTimeout()
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		baseURL+"/user/"+url.PathEscape("testuser2@example.com"),
		strings.NewReader(`{"blocked":true,"block_reason":"test block"}`),
	)
	require.NoError(t, err)
	req.Header.Add("Authorization", "testapikey123456789abcdefg")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	_, status := tests.ParseJSONBodyWithMessage(t, resp.Body)
	assert.Equal(t, "success", status)

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		baseURL+"/user/"+url.PathEscape("testuser2@example.com"),
		strings.NewReader(`{"blocked":false}`),
	)
	require.NoError(t, err)
	req.Header.Add("Authorization", "testapikey123456789abcdefg")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	_, status = tests.ParseJSONBodyWithMessage(t, resp.Body)
	assert.Equal(t, "success", status)
}
