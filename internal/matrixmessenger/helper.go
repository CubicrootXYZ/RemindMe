package matrixmessenger

import (
	"fmt"
	"regexp"
	"strings"
)

func makeResponse(newMsg, newMsgFormatted, respondMsg, respondMsgFormatted, respondTo, roomID, respondEventID string) (body, bodyFormatted string) {
	body = fmt.Sprintf("> <%s>%s\n\n%s", respondTo, respondMsg, newMsg)

	bodyFormatted = fmt.Sprintf("<mx-reply><blockquote><a href='https://matrix.to/#/%s/%s'>In reply to</a> <a href='https://matrix.to/#/%s'>%s</a><br />%s</blockquote></mx-reply>%s", roomID, respondEventID, respondTo, respondTo, respondMsgFormatted, newMsgFormatted)

	return body, bodyFormatted
}

func makeLinkToUser(userID string) (link string) {
	re := regexp.MustCompile("@(.+):")

	link = fmt.Sprintf("<a href=\"https://matrix.to/#/%s\">%s</a>", userID, re.Find([]byte(userID)))

	return
}

// StripReply removes the quoted reply from a message
func StripReply(msg string) string {
	strippedMsg := strings.Builder{}
	for _, line := range strings.Split(msg, "\n") {
		if strings.HasPrefix(line, ">") {
			continue
		}
		strippedMsg.WriteString(line)
	}

	return strippedMsg.String()
}

// StripReplyFormatted removes the quoted reply from a message
func StripReplyFormatted(msg string) string {
	re := regexp.MustCompile(`(?s)<mx-reply>.*<\/mx-reply>`)
	return re.ReplaceAllString(msg, "")
}
