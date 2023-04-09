package message

import (
	"regexp"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix"
)

var enableICalExportRegex = regexp.MustCompile("(?i)^(ical$|(show|give|list|send|write|)[ ]*(|me)[ ]*(the|)[ ]*(calendar|ical|cal|reminder|ics)[ ]+(link|url|uri|file))[ ]*$")

// EnableICalExportAction enables iCal in a channel.
type EnableICalExportAction struct {
	logger    gologger.Logger
	client    *mautrix.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
}

// Configure is called on startup and sets all dependencies.
func (action *EnableICalExportAction) Configure(
	logger gologger.Logger,
	client *mautrix.Client,
	messenger messenger.Messenger,
	matrixDB matrixdb.Service,
	db database.Service,
) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
}

// Name of the action.
func (action *EnableICalExportAction) Name() string {
	return "Enable iCal export"
}

// GetDocu returns the documentation for the action.
func (action *EnableICalExportAction) GetDocu() (title, explaination string, examples []string) {
	return "Enable iCal export",
		"Export events via a unique iCal URL.",
		[]string{"ical", "calendar link", "show me the calendar link"}
}

// Selector defines a regex on what messages the action should be used.
func (action *EnableICalExportAction) Selector() *regexp.Regexp {
	return enableICalExportRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *EnableICalExportAction) HandleEvent(event *matrix.MessageEvent) {
	// TODO
}
