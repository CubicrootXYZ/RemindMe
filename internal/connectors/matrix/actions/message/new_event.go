package message

import (
	"fmt"
	"regexp"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mapping"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

var DefaultEventTime = time.Minute
var ReminderRequestReactions = []string{"❌", "▶️", "⏩", "1️⃣", "4️⃣"}
var newEventActionRegex = regexp.MustCompile(".*")

// NewEventAction for new events. Should be the default message handler.
type NewEventAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
}

// Configure is called on startup and sets all dependencies.
func (action *NewEventAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
}

// Name of the action
func (action *NewEventAction) Name() string {
	return "New event"
}

func (action *NewEventAction) GetDocu() (title, explaination string, examples []string) {
	return "New event",
		"Add a new reminder.",
		[]string{"go shopping at monday", "buy milk at 5 pm", "ask boss for pay raise in 1 year"}
}

// Selector defines a regex on what messages the action should be used.
func (action *NewEventAction) Selector() *regexp.Regexp {
	return newEventActionRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *NewEventAction) HandleEvent(event *matrix.MessageEvent) {
	remindTime, err := format.ParseTime(event.Content.Body, event.Room.TimeZone, false)
	if err != nil {
		action.logger.Err(err)
		_ = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Sorry I was not able to understand the remind date and time from this message",
			event.Event.ID.String(),
			event.Content.Body,
			event.Event.Sender.String(),
			event.Room.RoomID,
		))
		return
	}

	dbEvent, err := action.db.NewEvent(&database.Event{
		Time:      remindTime,
		Duration:  DefaultEventTime,
		Message:   event.Content.Body,
		Active:    true,
		ChannelID: event.Channel.ID,
		InputID:   &event.Input.ID,
	})
	if err != nil {
		action.logger.Errorf("failed to save event to db: %v", err)
		return
	}

	message := mapping.MessageFromEvent(event)
	message.Type = matrixdb.MessageTypeNewEvent
	message.EventID = &dbEvent.ID
	_, err = action.matrixDB.NewMessage(message)
	if err != nil {
		action.logger.Errorf("failed to save message to db: %v", err)
		return
	}

	go func(evt *matrix.MessageEvent, dbEvent *database.Event) {
		msg := fmt.Sprintf("Successfully added new reminder (ID: %d) for %s", dbEvent.ID, format.ToLocalTime(dbEvent.Time, event.Room.TimeZone))

		response := messenger.PlainTextResponse(
			msg,
			evt.Event.ID.String(),
			evt.Content.Body,
			evt.Event.Sender.String(),
			evt.Room.RoomID,
		)

		resp, err := action.messenger.SendResponse(response)
		if err != nil {
			action.logger.Errorf("failed sending out message: %v", err)
			return
		}

		replyTo := event.Event.ID.String()
		message := mapping.MessageFromEvent(event)
		message.Type = matrixdb.MessageTypeNewEvent
		message.ID = resp.ExternalIdentifier
		message.ReplyToMessageID = &replyTo
		message.Incoming = false
		message.SendAt = time.Now()
		message.Body = msg
		message.BodyFormatted = msg
		message.EventID = &dbEvent.ID
		_, err = action.matrixDB.NewMessage(message)
		if err != nil {
			action.logger.Infof("Could not add message to database: %v", err)
		}
	}(event, dbEvent)

	for _, reaction := range ReminderRequestReactions {
		err = action.messenger.SendReactionAsync(&messenger.Reaction{
			Reaction:                  reaction,
			MessageExternalIdentifier: event.Event.ID.String(),
			ChannelExternalIdentifier: event.Room.RoomID,
		})
		if err != nil {
			action.logger.Err(err)
		}
	}
}
