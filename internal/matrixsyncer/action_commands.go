package matrixsyncer

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getActionCommands() *Action {
	action := &Action{
		Name:     "List all commands",
		Examples: []string{"show all commands", "list the commands", "commands"},
		Regex:    "(?i)(^(show|list)( all| the| my)( command| commands)$|commands|help)",
		Action:   s.actionCommands,
	}
	return action
}

// actionCommands lists all available commands
func (s *Syncer) actionCommands(evt *event.Event, channel *database.Channel) error {
	msg := formater.Formater{}

	msg.Title("Available Commands")
	msg.TextLine("You can interact with me in many ways, check out my features:")
	msg.NewLine()

	for _, action := range s.actions {
		msg.BoldLine(action.Name)

		if len(action.Examples) > 0 {
			msg.TextLine("Here are some examples how you can tell me to perform this action:")
			msg.List(action.Examples)
		}
		msg.NewLine()
	}

	msg.BoldLine("Make a new Reminder")
	msg.TextLine("To make a new reminder I will process all messages that are not part of one of the above commands. Try it with:")
	msg.List([]string{"Laundry at Sunday 12am", "Go shopping in 4 hours"})
	msg.NewLine()

	if len(s.reactionActions) > 0 {
		msg.SubTitle("Reactions")
		msg.TextLine("I am able to understand a few reactions you can give to a message.")
		msg.NewLine()

		for _, action := range s.reactionActions {
			msg.BoldLine(action.Name)
			msg.Text("Available for messages of the type " + string(action.Type) + ". Give the message one of these reactions: ")
			for _, reaction := range action.Keys {
				msg.Text(reaction + " ")
			}
			msg.NewLine()
		}
		msg.NewLine()
	}

	msg.SubTitle("Replies")
	msg.TextLine("I can also understand some of your replies to messages.")
	msg.NewLine()

	if len(s.replyActions) > 0 {
		for _, action := range s.replyActions {
			msg.BoldLine(action.Name)
			msg.TextLine("Answer to a message of the type " + string(action.ReplyToType) + " with: ")
			msg.List(action.Examples)
		}
	}

	msg.BoldLine("Change reminder time")
	msg.TextLine("You can achieve this with a reply to a message of the type REMINDER_REQUEST with one of this examples:")
	msg.List([]string{"sunday 5pm", "monday 15:57", "in 5 hours"})

	message, messageFormatted := msg.Build()

	_, err := s.messenger.SendFormattedMessage(message, messageFormatted, channel, database.MessageTypeActions, 0)
	return err
}
