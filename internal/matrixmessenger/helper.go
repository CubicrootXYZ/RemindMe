package matrixmessenger

import (
	"fmt"
	"regexp"
	"strings"
)

func makeResponse(newMsg, newMsgFormatted, respondMsg, respondMsgFormatted, respondTo, roomID, respondEventID string) (body, bodyFormatted string) {
	body = fmt.Sprintf("> <%s>%s\n\n%s", respondTo, respondMsg, newMsg)

	bodyFormatted = fmt.Sprintf("<mx-reply><blockquote><a href=\"https://matrix.to/#/%s/%s?via=%s\">In reply to</a> <a href=\"https://matrix.to/#/%s\">%s</a><br>%s</blockquote></mx-reply>%s", roomID, respondEventID, getHomeServerFromFullUsername(respondTo), respondTo, respondTo, respondMsgFormatted, newMsgFormatted)

	return body, bodyFormatted
}

func makeLinkToUser(userID string) (link string) {
	re := regexp.MustCompile("@(.+):")

	link = fmt.Sprintf(`<a href="https://matrix.to/#/%s">%s</a>`, userID, re.Find([]byte(userID)))

	return
}

func getHomeServerFromFullUsername(username string) string {
	if !strings.Contains(username, ":") {
		return "matrix.org"
	}

	return strings.Split(username, ":")[1]
}
