package message

import (
	"errors"
	"regexp"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mapping"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

var enableICalExportRegex = regexp.MustCompile("(?i)^(ical$|(show|give|list|send|write|)[ ]*(|me)[ ]*(the|)[ ]*(calendar|ical|cal|reminder|ics)[ ]+(link|url|uri|file))[ ]*$")

// EnableICalExportAction enables iCal in a channel.
type EnableICalExportAction struct {
	logger     gologger.Logger
	client     mautrixcl.Client
	messenger  messenger.Messenger
	matrixDB   matrixdb.Service
	db         database.Service
	icalBridge matrix.BridgeServiceICal
}

// Configure is called on startup and sets all dependencies.
func (action *EnableICalExportAction) Configure(
	logger gologger.Logger,
	client mautrixcl.Client,
	messenger messenger.Messenger,
	matrixDB matrixdb.Service,
	db database.Service,
	bridgeServices *matrix.BridgeServices,
) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
	action.icalBridge = bridgeServices.ICal
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
	icalOutput, err := action.getOrCreateIcalOutput(event)
	if err != nil {
		err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Whoopsie, that failed",
			event.Event.ID.String(),
			event.Content.Body,
			event.Event.Sender.String(),
			event.Room.RoomID,
		))
		if err != nil {
			action.logger.Err(err)
		}
		action.logger.Err(err)
		return
	}

	// Add message to database
	msg := mapping.MessageFromEvent(event)
	msg.Type = matrixdb.MessageTypeIcalExportEnable
	_, err = action.matrixDB.NewMessage(msg)
	if err != nil {
		action.logger.Err(err)
	}

	message := "Your calendar is ready 🥳: " + icalOutput.Token // TODO get baseurl
	resp, err := action.messenger.SendResponse(messenger.PlainTextResponse(
		message,
		event.Event.ID.String(),
		event.Content.Body,
		event.Event.Sender.String(),
		event.Room.RoomID,
	))
	if err != nil {
		action.logger.Err(err)
		return
	}

	msg = mapping.MessageFromEvent(event)
	msg.SendAt = resp.Timestamp
	msg.ID = resp.ExternalIdentifier
	msg.Incoming = false
	msg.Type = matrixdb.MessageTypeIcalExportEnable
	msg.Body = message
	msg.BodyFormatted = message
	_, err = action.matrixDB.NewMessage(msg)
	if err != nil {
		action.logger.Err(err)
	}
}

func (action *EnableICalExportAction) getOrCreateIcalOutput(event *matrix.MessageEvent) (*icaldb.IcalOutput, error) {
	for _, o := range event.Channel.Outputs {
		if o.OutputType == ical.OutputType {
			icalOutput, err := action.icalBridge.GetOutput(o.ID)
			if err == nil {
				return icalOutput, nil
			} else if !errors.Is(err, ical.ErrNotFound) {
				return nil, err
			}
		}
	}

	return action.icalBridge.NewOutput(event.Channel.ID)
}
