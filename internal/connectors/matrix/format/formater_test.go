package format

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
	shouldMsgHTML := "<h3>Title</h3><br><h4>Subtitle</h4><br><a href='domain.tld'>Link text</a><br><b>bold line of text</b><br><blockquote>quoted text</blockquote><br><br><br><i>italic line of text</i><br>Usual line of text<br>Just my text<br><ul><li>item1</li><li>item2</li></ul><br><span data-mx-spoiler>secret :D</span>"

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

func TestFormater_Bold(t *testing.T) {
	formater := Formater{}
	formater.Bold("This is bold")
	msg, msgFormatted := formater.Build()

	assert.Equal(t, "THIS IS BOLD", msg)
	assert.Equal(t, "<b>This is bold</b>", msgFormatted)
}

func TestFormater_Italic(t *testing.T) {
	formater := Formater{}
	formater.Italic("This is italic")
	msg, msgFormatted := formater.Build()

	assert.Equal(t, "This is italic", msg)
	assert.Equal(t, "<i>This is italic</i>", msgFormatted)
}
