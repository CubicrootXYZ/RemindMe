package formater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormater_StripReply(t *testing.T) {
	testCases := make(map[string]string)
	testCases["Hello"] = "> My name is fjfjfjf\nHello\n> I am cool"
	testCases["some random text "] = "some random text \n> I am cool"

	for should, msg := range testCases {
		is := StripReply(msg)
		assert.Equal(t, should, is)
	}
}

func TestFormater_StripReplyFormatted(t *testing.T) {
	testCases := make(map[string]string)
	testCases["Hello"] = "<mx-reply> My name is fjfjfjf</mx-reply>Hello<mx-reply> I am cool</mx-reply>"
	testCases["some random text "] = "some random text <mx-reply>> I am cool</mx-reply>"
	testCases["some random text "] = "<mx-reply>> I am cool</mx-reply>some random text "

	for should, msg := range testCases {
		is := StripReplyFormatted(msg)
		assert.Equal(t, should, is)
	}
}
