package tests

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// ParseJSONBodyWithMap parses the body of a JSON response.
// Returns the data field as map[string]interface{} and the status message as string.
func ParseJSONBodyWithMap(t *testing.T, body io.ReadCloser) (map[string]interface{}, string) {
	t.Helper()

	var bodyParsed map[string]interface{}
	err := json.NewDecoder(body).Decode(&bodyParsed)
	require.NoError(t, err)

	status, ok := bodyParsed["status"].(string)
	require.True(t, ok, "status is not a string")

	data, ok := bodyParsed["data"].(map[string]interface{})
	require.True(t, ok, "data field is malformed")

	return data, status
}

// ParseJSONBodyWithSlice parses the body of a JSON response.
// Returns the data field as map[string]interface{} and the status message as string.
func ParseJSONBodyWithSlice(t *testing.T, body io.ReadCloser) ([]interface{}, string) {
	t.Helper()

	var bodyParsed map[string]interface{}
	err := json.NewDecoder(body).Decode(&bodyParsed)
	require.NoError(t, err)

	status, ok := bodyParsed["status"].(string)
	require.True(t, ok, "status is not a string")

	data, ok := bodyParsed["data"].([]interface{})
	require.True(t, ok, "data field is malformed")

	return data, status
}
