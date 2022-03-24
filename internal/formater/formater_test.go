package formater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormater_Formater(t *testing.T) {
	formater := Formater{}
	formater.Title("Title")
	formater.SubTitle("Subtitle")
	formater.Link("Link text", "domain.tld")
	formater.NewLine()
	formater.BoldLine("bold line of text")
	formater.QuoteLine("quoted text")
	formater.DoubleNewLine()
	formater.ItalicLine("italic line of text")
	formater.TextLine("Usual line of text")
	formater.Text("Just my text")
	formater.NewLine()
	formater.List([]string{"item1", "item2"})
	formater.Spoiler("secret :D")

	msg, msgHTML := formater.Build()

	shouldMsg := "== TITLE ==\n==== SUBTITLE ====\ndomain.tld\nBOLD LINE OF TEXT\n> quoted text\n\n\nitalic line of text\nUsual line of text\nJust my text\n- item1\n- item2\nsecret :D"
	shouldMsgHTML := "<h3>Title</h3><h4>Subtitle</h4><a href='domain.tld'>Link text</a><br><b>bold line of text</b><br><blockquote>quoted text</blockquote><br><br><br><i>italic line of text</i><br>Usual line of text<br>Just my text<br><ul><li>item1</li><li>item2</li></ul><span data-mx-spoiler>secret :D</span>"

	assert.Equal(t, shouldMsg, msg)
	assert.Equal(t, shouldMsgHTML, msgHTML)
}

func TestFormater_Username(t *testing.T) {
	formater := Formater{}
	formater.Username("abcdefgh:matrix.org")
	msg, msgFormatted := formater.Build()

	assert.Equal(t, "abcdefgh", msg)
	assert.Equal(t, "<a href=\"https://matrix.to/#/@abcdefgh:matrix.org\">abcdefgh</a>", msgFormatted)
}

func TestFormater_UsernameWithOutDomain(t *testing.T) {
	formater := Formater{}
	formater.Username("abcdefgh")
	msg, msgFormatted := formater.Build()

	assert.Equal(t, "abcdefgh", msg)
	assert.Equal(t, "<a href=\"https://matrix.to/#/@abcdefgh\">abcdefgh</a>", msgFormatted)
}

func TestFormater_UsernameWithAt(t *testing.T) {
	formater := Formater{}
	formater.Username("@abcdefgh:matrix.org")
	msg, msgFormatted := formater.Build()

	assert.Equal(t, "abcdefgh", msg)
	assert.Equal(t, "<a href=\"https://matrix.to/#/@abcdefgh:matrix.org\">abcdefgh</a>", msgFormatted)
}
