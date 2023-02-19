package format

import (
	"regexp"
	"strings"
)

// FullUsername assembles the full username from username and homerserver.
// Username can already be a full username, it will be returned without change.
func FullUsername(username string, homeserver string) string {
	if !strings.HasPrefix(username, "@") {
		username = "@" + username
	}

	if strings.Contains(username, ":") {
		return username
	}

	return username + ":" + strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(homeserver, "http://"), "https://"), "/")
}

// GetUsernameFromLink searches for an user link in the given input and extracts the username.
// Returns an empty string if no user is found.
func GetUsernameFromLink(link string) string {
	r := regexp.MustCompile(`https:\/\/matrix.to\/#\/[^"'>]+`)

	url := r.Find([]byte(link))
	if url == nil {
		return ""
	}

	return strings.TrimPrefix(string(url), "https://matrix.to/#/")
}
