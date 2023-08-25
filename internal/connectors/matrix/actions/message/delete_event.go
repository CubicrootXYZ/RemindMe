package message

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mapping"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/msghelper"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

var deleteEventActionRegex = regexp.MustCompile("(?i)(^(delete|remove)[ ]*(reminder|)[ ]+[0-9]+)[ ]*$")

// DeleteEventAction acts as a template for new actions.
type DeleteEventAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
	storer    msghelper.Storer
}

// Configure is called on startup and sets all dependencies.
func (action *DeleteEventAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
	action.storer = *msghelper.NewStorer(matrixDB, messenger, logger)
}

// Name of the action
func (action *DeleteEventAction) Name() string {
	return "Delete Event"
}

// GetDocu returns the documentation for the action.
func (action *DeleteEventAction) GetDocu() (title, explaination string, examples []string) {
	return "Delete event",
		"Delete an event by its ID.",
		[]string{"delete reminder 1", "remove 68"}
}

// Selector defines a regex on what messages the action should be used.
func (action *DeleteEventAction) Selector() *regexp.Regexp {
	return deleteEventActionRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *DeleteEventAction) HandleEvent(event *matrix.MessageEvent) {
	id, err := getIDFromSentence(event.Content.Body)
	if err != nil {
		action.logger.Err(err)
		err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Ups, can not find an ID in there.",
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

	events, err := action.db.ListEvents(&database.ListEventsOpts{
		IDs:       []uint{uint(id)},
		ChannelID: &event.Channel.ID,
	})
	if err != nil {
		action.logger.Err(err)
		err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Sorry, an error appeared.",
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

	if len(events) != 1 {
		err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"I could not find that event in my database.",
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

	dbMessage := mapping.MessageFromEvent(event)
	dbMessage.Type = matrixdb.MessageTypeEventDelete
	dbMessage.EventID = &events[0].ID
	_, err = action.matrixDB.NewMessage(dbMessage)
	if err != nil {
		action.logger.Err(err)
		return
	}

	err = action.db.DeleteEvent(&events[0])
	if err != nil {
		action.logger.Err(err)
		err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Sorry, an error appeared.",
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

	go action.storer.SendAndStoreResponse("Deleted event \""+events[0].Message+"\"", matrixdb.MessageTypeEventDelete, *event, msghelper.WithEventID(events[0].ID))
}

func getIDFromSentence(value string) (int, error) {
	splitUp := strings.Split(value, " ")
	if len(splitUp) == 0 {
		return 0, errors.New("empty string does not contain integer")
	}

	integerString := splitUp[len(splitUp)-1]

	return strconv.Atoi(integerString)
}
