package message

import (
	"errors"
	"regexp"

	"github.com/CubicrootXYZ/gologger"
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

var regenIcalTokenActionRegex = regexp.MustCompile("(?i)^(make|generate|)[ ]*(renew|generate|delete|regenerate|renew|new)[ ]*(the|a|)[ ]+(ical|calendar|token|secret)[ ]*(token|secret|)[ ]*$")

// RegenICalTokenAction enables iCal in a channel.
type RegenICalTokenAction struct {
	logger     gologger.Logger
	client     mautrixcl.Client
	messenger  messenger.Messenger
	matrixDB   matrixdb.Service
	db         database.Service
	icalBridge matrix.BridgeServiceICal
}

// Configure is called on startup and sets all dependencies.
func (action *RegenICalTokenAction) Configure(
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
func (action *RegenICalTokenAction) Name() string {
	return "Regenerate iCal export token"
}

// GetDocu returns the documentation for the action.
func (action *RegenICalTokenAction) GetDocu() (title, explaination string, examples []string) {
	return "Regenerate iCal export token",
		"Generate a new token for the iCal export. Previously generated URLs won't work anymore.",
		[]string{"renew the calendar secret", "generate token"}
}

// Selector defines a regex on what messages the action should be used.
func (action *RegenICalTokenAction) Selector() *regexp.Regexp {
	return regenIcalTokenActionRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *RegenICalTokenAction) HandleEvent(event *matrix.MessageEvent) {
	_, url, err := action.getIcalOutput(event)
	if err != nil {
		msg := "Oh no, this did not work 😰"
		if errors.Is(err, ical.ErrNotFound) {
			msg = "It looks like iCal output is not set up for this channel. Set it up first."
		} else {
			action.logger.Err(err)
		}
		err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			msg,
			event.Event.ID.String(),
			event.Content.Body,
			event.Event.Sender.String(),
			event.Room.RoomID,
		))
		if err != nil {
			action.logger.Err(err)
		}
		return
	}

	// Add message to database
	msg := mapping.MessageFromEvent(event)
	msg.Type = matrixdb.MessageTypeIcalRegenToken
	_, err = action.matrixDB.NewMessage(msg)
	if err != nil {
		action.logger.Err(err)
	}

	msgBuilder := format.Formater{}
	msgBuilder.Text("Your new secret calendar URL is: ")
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
		action.logger.Err(err)
		return
	}

	msg = mapping.MessageFromEvent(event)
	msg.SendAt = resp.Timestamp
	msg.ID = resp.ExternalIdentifier
	msg.Incoming = false
	msg.Type = matrixdb.MessageTypeIcalRegenToken
	msg.Body = response
	msg.BodyFormatted = response
	_, err = action.matrixDB.NewMessage(msg)
	if err != nil {
		action.logger.Err(err)
	}
}

func (action *RegenICalTokenAction) getIcalOutput(event *matrix.MessageEvent) (*icaldb.IcalOutput, string, error) {
	for _, o := range event.Channel.Outputs {
		if o.OutputType == ical.OutputType {
			icalOutput, url, err := action.icalBridge.GetOutput(o.OutputID, true)
			if err == nil {
				return icalOutput, url, nil
			} else if !errors.Is(err, ical.ErrNotFound) {
				return nil, "", err
			}
		}
	}

	return nil, "", ical.ErrNotFound
}
