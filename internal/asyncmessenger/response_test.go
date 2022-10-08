package asyncmessenger

import (
	"testing"

	"github.com/tj/assert"
)

func TestResponse_GetResponseMessage(t *testing.T) {
	response := Response{
		"My Answer",
		"My <b>Answer</b>",
		"The message",
		"The message",
		"@user:example.com",
		"!abcde-1234",
		"!12345",
	}

	msg, msgFormatted := response.getResponseMessage()

	assert.Equal(t, "> <@user:example.com>The message\n\nMy Answer", msg)
	assert.Equal(t, "<mx-reply><blockquote><a href=\"https://matrix.to/#/!12345/!abcde-1234?via=example.com\">In reply to</a> <a href=\"https://matrix.to/#/@user:example.com\">@user:example.com</a><br>The message</blockquote></mx-reply>My <b>Answer</b>", msgFormatted)
}