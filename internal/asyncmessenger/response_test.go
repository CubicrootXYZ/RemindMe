package asyncmessenger

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	msg, msgFormatted := response.GetResponseMessage()

	assert.Equal(t, "> <@user:example.com>The message\n\nMy Answer", msg)
	assert.Equal(t, "<mx-reply><blockquote><a href=\"https://matrix.to/#/!12345/!abcde-1234?via=example.com\">In reply to</a> <a href=\"https://matrix.to/#/@user:example.com\">@user:example.com</a><br>The message</blockquote></mx-reply>My <b>Answer</b>", msgFormatted)
}

func TestResponse_GetResponseMessage_DoubleResponse(t *testing.T) {
	response := Response{
		"My Answer",
		"My <b>Answer</b>",
		"> <@user:example.com>Original message\n\nThe message",
		"<mx-reply><blockquote><a href=\"https://matrix.to/#/!12345/!abcde-1234?via=example.com\">In reply to</a> <a href=\"https://matrix.to/#/@user:example.com\">@user:example.com</a><br>Original message</blockquote></mx-reply>The message",
		"@user:example.com",
		"!abcde-1234",
		"!12345",
	}

	msg, msgFormatted := response.GetResponseMessage()

	assert.Equal(t, "> <@user:example.com>The message\n\nMy Answer", msg)
	assert.Equal(t, "<mx-reply><blockquote><a href=\"https://matrix.to/#/!12345/!abcde-1234?via=example.com\">In reply to</a> <a href=\"https://matrix.to/#/@user:example.com\">@user:example.com</a><br>The message</blockquote></mx-reply>My <b>Answer</b>", msgFormatted)
}
