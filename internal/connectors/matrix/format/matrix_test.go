package format_test

import (
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/stretchr/testify/assert"
)

func TestGetMatrixLinkForUser(t *testing.T) {
	result := format.GetMatrixLinkForUser("@mybot:example.com")
	assert.Equal(t, `<a href="https://matrix.to/#/@mybot:example.com">@mybot</a>`, result)
}

func TestGetMatrixLinkForUser_SpecialCharacters(t *testing.T) {
	result := format.GetMatrixLinkForUser("@myböté:exämple.com")
	assert.Equal(t, `<a href="https://matrix.to/#/@myböté:exämple.com">@myböté</a>`, result)
}

func TestGetMatrixLinkForUser_MissingInstance(t *testing.T) {
	result := format.GetMatrixLinkForUser("@mybot")
	assert.Equal(t, `<a href="https://matrix.to/#/@mybot">@mybot</a>`, result)
}

func TestGetHomeserverFromUserID(t *testing.T) {
	result := format.GetHomeserverFromUserID("@user:example.com")
	assert.Equal(t, `example.com`, result)
}

func TestGetHomeserverFromUserID_SpecialCharacters(t *testing.T) {
	result := format.GetHomeserverFromUserID("@user:éxämple.com")
	assert.Equal(t, `éxämple.com`, result)
}

func TestGetHomeserverFromUserID_EmptyHomeserver(t *testing.T) {
	result := format.GetHomeserverFromUserID("@user:")
	assert.Empty(t, result)
}

func TestGetHomeserverFromUserID_NoHomeserver(t *testing.T) {
	result := format.GetHomeserverFromUserID("@user")
	assert.Equal(t, `matrix.org`, result)
}

func TestGetUsernameFromUserIdentifier(t *testing.T) {
	result := format.GetUsernameFromUserIdentifier("@testuser:example.org")
	assert.Equal(t, "testuser", result)
}

func TestGetUsernameFromUserIdentifier_EmptyHomeserver(t *testing.T) {
	result := format.GetUsernameFromUserIdentifier("@testuser:")
	assert.Equal(t, "testuser", result)
}

func TestGetUsernameFromUserIdentifier_NoHomeserver(t *testing.T) {
	result := format.GetUsernameFromUserIdentifier("@testuser")
	assert.Equal(t, "testuser", result)
}

func TestGetUsernameFromUserIdentifier_InvalidString(t *testing.T) {
	result := format.GetUsernameFromUserIdentifier("testuser")
	assert.Equal(t, "testuser", result)
}
