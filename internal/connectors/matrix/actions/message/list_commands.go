package message

import (
	"regexp"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/reaction"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/reply"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/msghelper"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

var listCommandsRegex = regexp.MustCompile("(?i)^(((show|list)( all| the| my)( command| commands))|commands|help)[ ]*$")

// ListCommandsAction sets the time for the daily reminder.
type ListCommandsAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
	storer    *msghelper.Storer
}

// Configure is called on startup and sets all dependencies.
func (action *ListCommandsAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
	action.storer = msghelper.NewStorer(matrixDB, messenger, logger)
}

// Name of the action
func (action *ListCommandsAction) Name() string {
	return "List Commands"
}

// GetDocu returns the documentation for the action.
func (action *ListCommandsAction) GetDocu() (title, explaination string, examples []string) {
	return "List Commands",
		"List available commands",
		[]string{"show all commands", "list the commands", "commands"}
}

// Selector defines a regex on what messages the action should be used.
func (action *ListCommandsAction) Selector() *regexp.Regexp {
	return listCommandsRegex
}

// HandleEvent is where the message event get's send to if it matches the Selector.
func (action *ListCommandsAction) HandleEvent(event *matrix.MessageEvent) {
	action.listMessageCommands(event)
	action.listReplyCommands(event)
	action.listReactions(event)
}

func (action *ListCommandsAction) listMessageCommands(event *matrix.MessageEvent) {
	msg := format.Formater{}
	msg.Title("RemindMe Commands")
	msg.TextLine("Try messaging me with the following messages and I will give my best to assist you!")

	for _, action := range []interface {
		GetDocu() (string, string, []string)
	}{
		&AddUserAction{},
		&ChangeEventAction{},
		&ChangeTimezoneAction{},
		&DeleteEventAction{},
		&EnableICalExportAction{},
		&ListCommandsAction{},
		&ListEventsAction{},
		&NewEventAction{},
		&RegenICalTokenAction{},
		&SetDailyReminderAction{},
	} {
		title, explain, examples := action.GetDocu()
		msg.BoldLine(title)
		msg.TextLine(explain)
		msg.TextLine("For example: ")
		msg.List(examples)
	}

	message, messageFormatted := msg.Build()
	go action.storer.SendAndStoreMessage(
		message,
		messageFormatted,
		matrixdb.MessageTypeListCommands,
		*event,
	)
}

func (action *ListCommandsAction) listReplyCommands(event *matrix.MessageEvent) {
	msg := format.Formater{}
	msg.Title("RemindMe Reply Commands")
	msg.TextLine("I can understand the following messages if you reply with them to certain messages.")

	for _, action := range []interface {
		GetDocu() (string, string, []string)
	}{
		&reply.ChangeTimeAction{},
		&reply.DeleteEventAction{},
		&reply.MakeRecurringAction{},
	} {
		title, explain, examples := action.GetDocu()
		msg.BoldLine(title)
		msg.TextLine(explain)
		msg.TextLine("For example: ")
		msg.List(examples)
	}

	message, messageFormatted := msg.Build()
	go action.storer.SendAndStoreMessage(
		message,
		messageFormatted,
		matrixdb.MessageTypeListCommands,
		*event,
	)
}

func (action *ListCommandsAction) listReactions(event *matrix.MessageEvent) {
	msg := format.Formater{}
	msg.Title("RemindMe Reactions")
	msg.TextLine("You can react to some messages to let me know I should act.")

	for _, action := range []interface {
		GetDocu() (string, string, []string)
	}{
		&reaction.AddTimeAction{},
		&reaction.DeleteEventAction{},
		&reaction.MarkDoneAction{},
	} {
		title, explain, _ := action.GetDocu()
		msg.BoldLine(title)
		msg.TextLine(explain)
	}

	message, messageFormatted := msg.Build()
	go action.storer.SendAndStoreMessage(
		message,
		messageFormatted,
		matrixdb.MessageTypeListCommands,
		*event,
	)
}
