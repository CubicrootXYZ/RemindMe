package format

import (
	"regexp"
	"strings"
)

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
	re := regexp.MustCompile(`(?s)<mx-reply>.*?<\/mx-reply>`)
	return re.ReplaceAllString(msg, "")
}
