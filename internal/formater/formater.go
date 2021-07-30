package formater

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

// Quote quotes the text
func (f *Formater) Quote(text string) {
	f.msg.WriteString("> " + text + "\n")
	f.msgFormatted.WriteString("<blockquote>" + text + "</blockquote>" + "<br>")
}

// ItalicLine adds italic text
func (f *Formater) ItalicLine(text string) {
	f.msg.WriteString(text + "\n")
	f.msgFormatted.WriteString("<i>" + text + "</i><br>")
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

// Build returns the build formatted and unformatted messages
func (f *Formater) Build() (message, messageFormatted string) {
	return f.msg.String(), f.msgFormatted.String()
}
