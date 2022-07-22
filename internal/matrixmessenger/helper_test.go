package matrixmessenger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeResponseOnSuccess(t *testing.T) {
	response, responseFormatted := makeResponse("testreply", "testreply", "testmessage", "testmessage", "@cubetest:matrix.org", "!ABCDEFGHIJ:matrix.org", "$1234ABCDE")

	assert.Equal(t, "<mx-reply><blockquote><a href=\"https://matrix.to/#/!ABCDEFGHIJ:matrix.org/$1234ABCDE?via=matrix.org\">In reply to</a> <a href=\"https://matrix.to/#/@cubetest:matrix.org\">@cubetest:matrix.org</a><br>testmessage</blockquote></mx-reply>testreply", responseFormatted)
	assert.Equal(t, "> <@cubetest:matrix.org>testmessage\n\ntestreply", response)
}
