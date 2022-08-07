package matrixsyncer

import (
	"regexp"
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

func (s *Syncer) getActionCommands() *types.Action {
	action := &types.Action{
		Name:     "List all commands",
		Examples: []string{"show all commands", "list the commands", "commands"},
		Regex:    regexp.MustCompile("(?i)(^(show|list)( all| the| my)( command| commands)$|commands|help)"),
		Action:   s.actionCommands,
	}
	return action
}

// actionCommands lists all available commands
func (s *Syncer) actionCommands(evt *types.MessageEvent, channel *database.Channel) error {
	msg := formater.Formater{}

	msg.Title("Available Commands")
	msg.TextLine("You can interact with me in many ways, check out my features:")
	msg.NewLine()

	for _, action := range s.getActions() {
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

	reactionActions := s.getReactionActions()
	if len(reactionActions) > 0 {
		msg.SubTitle("Reactions")
		msg.TextLine("I am able to understand a few reactions you can give to a message.")
		msg.NewLine()

		nameToActionToType := make(map[string]map[string][]string)
		for _, action := range reactionActions {
			if _, exists := nameToActionToType[action.Name]; !exists {
				nameToActionToType[action.Name] = make(map[string][]string)
			}

			for _, key := range action.Keys {
				if _, exists := nameToActionToType[action.Name][key]; !exists {
					nameToActionToType[action.Name][key] = make([]string, 0)
				}

				nameToActionToType[action.Name][key] = append(nameToActionToType[action.Name][key], string(action.Type))
			}
		}

		for actionName, keyToActionType := range nameToActionToType {
			msg.Bold(actionName + ": ")

			keys := make([]string, 0)
			actionTypes := make([]string, 0)
			for key, actionType := range keyToActionType {
				keys = append(keys, key)
				actionTypes = append(actionTypes, actionType...)
			}

			msg.Text(strings.Join(keys, ", "))
			msg.Text(" avalaible on " + strings.Join(actionTypes, ", "))
		}

		msg.NewLine()
	}

	msg.SubTitle("Replies")
	msg.TextLine("I can also understand some of your replies to messages.")
	msg.NewLine()

	replyActions := s.getReplyActions()
	if len(replyActions) > 0 {
		for _, action := range replyActions {
			msg.BoldLine(action.Name)
			msg.Text("Answer to a message of the type ")
			if formater.EqMessageType(action.ReplyToTypes, database.MessageTypesWithReminder) {
				msg.Text("reminder or reminder edits")
			} else {
				for _, rtt := range action.ReplyToTypes {
					msg.Text(string(rtt) + " ")
				}
			}

			msg.TextLine(" with: ")
			msg.List(action.Examples)
		}
	}

	msg.BoldLine("Change reminder time")
	msg.TextLine("You can achieve this with a reply to a message of the type reminder or reminder edits with one of this examples:")
	msg.List([]string{"sunday 5pm", "monday 15:57", "in 5 hours"})

	message, messageFormatted := msg.Build()

	_, err := s.messenger.SendFormattedMessage(message, messageFormatted, channel, database.MessageTypeActions, 0)
	return err
}
