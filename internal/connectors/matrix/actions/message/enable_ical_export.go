package message

import (
	"errors"
	"log/slog"
	"regexp"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mapping"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

var enableICalExportActionRegex = regexp.MustCompile("(?i)^(ical$|(show|give|list|send|write|)[ ]*(|me)[ ]*(the|)[ ]*(calendar|ical|cal|reminder|ics)[ ]+(link|url|uri|file))[ ]*$")

// EnableICalExportAction enables iCal in a channel.
type EnableICalExportAction struct {
	logger     *slog.Logger
	client     mautrixcl.Client
	messenger  messenger.Messenger
	matrixDB   matrixdb.Service
	db         database.Service
	icalBridge matrix.BridgeServiceICal
}

// Configure is called on startup and sets all dependencies.
func (action *EnableICalExportAction) Configure(
	logger *slog.Logger,
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
	return enableICalExportActionRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *EnableICalExportAction) HandleEvent(event *matrix.MessageEvent) {
	_, url, err := action.getOrCreateIcalOutput(event)
	if err != nil {
		err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Whoopsie, that failed",
			event.Event.ID.String(),
			event.Content.Body,
			event.Event.Sender.String(),
			event.Room.RoomID,
		))
		if err != nil {
			action.logger.Error("failed to send response", "error", err)
		}
		action.logger.Error("failed to get/create iCal output", "error", err)
		return
	}

	// Add message to database
	msg := mapping.MessageFromEvent(event)
	msg.Type = matrixdb.MessageTypeIcalExportEnable
	_, err = action.matrixDB.NewMessage(msg)
	if err != nil {
		action.logger.Error("failed to save message to database", "error", err)
	}

	msgBuilder := format.Formater{}
	msgBuilder.Text("Your calendar is ready 🥳: ")
	msgBuilder.Link(url, url)

	response, _ := msgBuilder.Build()

	resp, err := action.messenger.SendResponse(messenger.PlainTextResponse(
		response,
		event.Event.ID.String(),
		event.Content.Body,
		event.Event.Sender.String(),
		event.Room.RoomID,
	))
	if err != nil {
		action.logger.Error("failed to send response", "error", err)
		return
	}

	msg = mapping.MessageFromEvent(event)
	msg.SendAt = resp.Timestamp
	msg.ID = resp.ExternalIdentifier
	msg.Incoming = false
	msg.Type = matrixdb.MessageTypeIcalExportEnable
	msg.Body = response
	msg.BodyFormatted = response
	_, err = action.matrixDB.NewMessage(msg)
	if err != nil {
		action.logger.Error("failed to save response to database", "error", err)
	}
}

func (action *EnableICalExportAction) getOrCreateIcalOutput(event *matrix.MessageEvent) (*icaldb.IcalOutput, string, error) {
	for _, o := range event.Channel.Outputs {
		if o.OutputType == ical.OutputType {
			icalOutput, url, err := action.icalBridge.GetOutput(o.OutputID, false)
			if err == nil {
				return icalOutput, url, nil
			} else if !errors.Is(err, ical.ErrNotFound) {
				return nil, "", err
			}
		}
	}

	return action.icalBridge.NewOutput(event.Channel.ID)
}
