package format

import "strings"

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
