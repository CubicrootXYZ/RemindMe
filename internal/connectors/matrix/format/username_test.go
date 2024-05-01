package format_test

import (
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/stretchr/testify/assert"
)

func TestFullUsername(t *testing.T) {
	type testCase struct {
		name                 string
		username             string
		homeserver           string
		expectedFullUsername string
	}

	testCases := []testCase{
		{
			name:                 "all empty",
			expectedFullUsername: "@:",
		},
		{
			name:                 "only username",
			username:             "user",
			expectedFullUsername: "@user:",
		},
		{
			name:                 "only homeserver",
			homeserver:           "example.com",
			expectedFullUsername: "@:example.com",
		},
		{
			name:                 "slim username, http homeserver",
			username:             "user",
			homeserver:           "http://example.com",
			expectedFullUsername: "@user:example.com",
		},
		{
			name:                 "slim username, https homeserver",
			username:             "user",
			homeserver:           "https://example.com",
			expectedFullUsername: "@user:example.com",
		},
		{
			name:                 "full username, https homeserver",
			username:             "@user:example.com",
			homeserver:           "https://example.com",
			expectedFullUsername: "@user:example.com",
		},
		{
			name:                 "full username without @, https homeserver",
			username:             "user:example.com",
			homeserver:           "https://example.com",
			expectedFullUsername: "@user:example.com",
		},
		{
			name:                 "full username, no homeserver",
			username:             "@user:example.com",
			expectedFullUsername: "@user:example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualFullUsername := format.FullUsername(tc.username, tc.homeserver)

			assert.Equal(t, tc.expectedFullUsername, actualFullUsername)
		})
	}
}

func TestGetUserFromLink(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name: "empty string",
		},
		{
			name:  "alphanumeric string",
			input: "abcde1234",
		},
		{
			name:           "alphanumeric link",
			input:          "https://matrix.to/#/username1",
			expectedOutput: "username1",
		},
		{
			name:           "full link",
			input:          "https://matrix.to/#/@user:example.com",
			expectedOutput: "@user:example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedOutput, format.GetUsernameFromLink(tc.input))
		})
	}
}
