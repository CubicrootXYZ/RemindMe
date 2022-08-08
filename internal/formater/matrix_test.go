package formater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMatrixLinkForUser(t *testing.T) {
	result := GetMatrixLinkForUser("@mybot:example.com")
	assert.Equal(t, `<a href="https://matrix.to/#/@mybot:example.com">@mybot</a>`, result)
}

func TestGetMatrixLinkForUser_SpecialCharacters(t *testing.T) {
	result := GetMatrixLinkForUser("@myböté:exämple.com")
	assert.Equal(t, `<a href="https://matrix.to/#/@myböté:exämple.com">@myböté</a>`, result)
}

func TestGetMatrixLinkForUser_MissingInstance(t *testing.T) {
	result := GetMatrixLinkForUser("@mybot")
	assert.Equal(t, `<a href="https://matrix.to/#/@mybot">@mybot</a>`, result)
}

func TestGetHomeserverFromUserID(t *testing.T) {
	result := GetHomeserverFromUserID("@user:example.com")
	assert.Equal(t, `example.com`, result)
}

func TestGetHomeserverFromUserID_SpecialCharacters(t *testing.T) {
	result := GetHomeserverFromUserID("@user:éxämple.com")
	assert.Equal(t, `éxämple.com`, result)
}

func TestGetHomeserverFromUserID_EmptyHomeserver(t *testing.T) {
	result := GetHomeserverFromUserID("@user:")
	assert.Equal(t, ``, result)
}

func TestGetHomeserverFromUserID_NoHomeserver(t *testing.T) {
	result := GetHomeserverFromUserID("@user")
	assert.Equal(t, `matrix.org`, result)
}
