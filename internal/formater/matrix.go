package formater

import (
	"fmt"
	"regexp"
	"strings"
)

var regexUsernameFromUserID = regexp.MustCompile("@([^:]+)")

// GetMatrixLinkForUser creates a clickable link pointing to the given user id
func GetMatrixLinkForUser(userID string) string {
	link := fmt.Sprintf(`<a href="https://matrix.to/#/%s">%s</a>`, userID, regexUsernameFromUserID.Find([]byte(userID)))

	return link
}

// GetHomeserverFromUserID returns the homeserver from a user id
func GetHomeserverFromUserID(userID string) string {
	if !strings.Contains(userID, ":") {
		return "matrix.org"
	}

	return strings.Split(userID, ":")[1]
}
