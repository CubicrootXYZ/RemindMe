package format

import "strings"

// Formater is a string builder with additional formating options
type Formater struct {
	msg          strings.Builder
	msgFormatted strings.Builder
}

// Title adds a title to the formater
func (f *Formater) Title(title string) {
	f.msg.WriteString("== " + strings.ToUpper(title) + " ==\n")
	f.msgFormatted.WriteString("<h3>" + title + "</h3>")
}

// SubTitle adds a subtitle to the formater
func (f *Formater) SubTitle(title string) {
	f.msg.WriteString("==== " + strings.ToUpper(title) + " ====\n")
	f.msgFormatted.WriteString("<h4>" + title + "</h4>")
}

// Link adds a link to the formater
func (f *Formater) Link(linkText string, url string) {
	f.msg.WriteString(url)
	f.msgFormatted.WriteString("<a href='" + url + "'>" + linkText + "</a>")
}

// NewLine adds a new line
func (f *Formater) NewLine() {
	f.msg.WriteString("\n")
	f.msgFormatted.WriteString("<br>")
}

// DoubleNewLine adds two new lines
func (f *Formater) DoubleNewLine() {
	f.msg.WriteString("\n\n")
	f.msgFormatted.WriteString("<br><br>")
}

// BoldLine adds bold text
func (f *Formater) BoldLine(text string) {
	f.msg.WriteString(strings.ToUpper(text) + "\n")
	f.msgFormatted.WriteString("<b>" + text + "</b><br>")
}

// Bold adds bold text
func (f *Formater) Bold(text string) {
	f.msg.WriteString(strings.ToUpper(text))
	f.msgFormatted.WriteString("<b>" + text + "</b>")
}

// QuoteLine quotes the text
func (f *Formater) QuoteLine(text string) {
	f.msg.WriteString("> " + text + "\n")
	f.msgFormatted.WriteString("<blockquote>" + text + "</blockquote>" + "<br>")
}

// ItalicLine adds italic text
func (f *Formater) ItalicLine(text string) {
	f.msg.WriteString(text + "\n")
	f.msgFormatted.WriteString("<i>" + text + "</i><br>")
}

// Italic adds italic text
func (f *Formater) Italic(text string) {
	f.msg.WriteString(text)
	f.msgFormatted.WriteString("<i>" + text + "</i>")
}

// TextLine adds the text
func (f *Formater) TextLine(text string) {
	f.msg.WriteString(text + "\n")
	f.msgFormatted.WriteString(text + "<br>")
}

// Text adds the text
func (f *Formater) Text(text string) {
	f.msg.WriteString(text)
	f.msgFormatted.WriteString(text)
}

// List adds a list to the message
func (f *Formater) List(items []string) {
	f.msgFormatted.WriteString("<ul>")
	for _, item := range items {
		f.msg.WriteString("- " + item + "\n")
		f.msgFormatted.WriteString("<li>" + item + "</li>")
	}
	f.msgFormatted.WriteString("</ul>")
}

// Spoiler adds a spoiler to the message
func (f *Formater) Spoiler(text string) {
	f.msg.WriteString(text)
	f.msgFormatted.WriteString("<span data-mx-spoiler>")
	f.msgFormatted.WriteString(text)
	f.msgFormatted.WriteString("</span>")
}

// Username adds a username reference to the message
func (f *Formater) Username(username string) {
	if len(username) == 0 {
		return
	}

	if username[0:1] != "@" {
		username = "@" + username
	}

	f.msgFormatted.WriteString("<a href=\"https://matrix.to/#/")
	f.msgFormatted.WriteString(username)
	f.msgFormatted.WriteString("\">")
	f.msgFormatted.WriteString(strings.TrimPrefix(strings.Split(username, ":")[0], "@"))
	f.msgFormatted.WriteString("</a>")

	f.msg.WriteString(strings.TrimPrefix(strings.Split(username, ":")[0], "@"))
}

// Build returns the build formatted and unformatted messages
func (f *Formater) Build() (message, messageFormatted string) {
	return f.msg.String(), f.msgFormatted.String()
}
