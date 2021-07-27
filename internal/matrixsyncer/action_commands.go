package matrixsyncer

import (
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getActionCommands() *Action {
	action := &Action{
		Name:     "List all commands",
		Examples: []string{"show all commands", "list the commands", "commands"},
		Regex:    "(?i)(^(show|list)( all| the| my)( command| commands)$|commands)",
		Action:   s.actionCommands,
	}
	return action
}

// actionCommands lists all available commands
func (s *Syncer) actionCommands(evt *event.Event, channel *database.Channel) error {
	msg := strings.Builder{}
	msgFormatted := strings.Builder{}

	msg.WriteString("= Available Commands = \nYou can interact with me in many ways, check out my features: \n\n")
	msgFormatted.WriteString("<h3>Available Commands</h3><You can interact with me in many ways, check out my features: <br><br>")

	for _, action := range s.actions {
		msg.WriteString("== " + action.Name + " ==\n")
		msgFormatted.WriteString("<b>" + action.Name + "</b><br>")

		if len(action.Examples) > 0 {
			msg.WriteString("Here are some examples how you can tell me to perform this action:\n")
			msgFormatted.WriteString("Here are some examples how you can tell me to perform this action:<br>")
			for _, example := range action.Examples {
				msg.WriteString(example + "\n")
				msgFormatted.WriteString("<i>" + example + "</i><br>")
			}
		}

		msg.WriteString("\n")
		msgFormatted.WriteString("<br>")

	}

	msg.WriteString("== Make a new Reminder ==\nTo make a new reminder I will process all messages that are not part of one of the above commands. Try it with:\nLaundry at Sunday 12am\nGo shopping in 4 hours")
	msgFormatted.WriteString("<b>Make a new Reminder </b><br>To make a new reminder I will process all messages that are not part of one of the above commands. Try it with:<br><i>Laundry at Sunday 12am</i><br><i>Go shopping in 4 hours</i>")

	_, err := s.messenger.SendFormattedMessage(msg.String(), msgFormatted.String(), channel.ChannelIdentifier)
	return err
}
