package matrixsyncer

import (
	"fmt"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getReactionActionDelete(rat ReactionActionType) *ReactionAction {
	action := &ReactionAction{
		Name:   "Delete a reminder",
		Keys:   []string{"❌"},
		Action: s.reactionActionDeleteReminder,
		Type:   rat,
	}
	return action
}

func (s *Syncer) reactionActionDeleteReminder(message *database.Message, content *event.ReactionEventContent, evt *event.Event, channel *database.Channel) error {
	var msg string
	var msgFormatted string
	reminder, err := s.daemon.Database.DeleteReminder(message.ReminderID)
	if err != nil {
		msg = fmt.Sprintf("Sorry, I could not delete the reminder %d.", message.ReminderID)
		msgFormatted = msg
		_, err = s.messenger.SendFormattedMessage(msg, msgFormatted, channel.ChannelIdentifier)
		return err
	}

	err = s.messenger.DeleteMessage(message.ExternalIdentifier, channel.ChannelIdentifier)
	if err != nil {
		return err
	}

	msg = fmt.Sprintf("I deleted the reminder \"%s\" (at %s) for you.", reminder.Message, reminder.RemindTime.Format("15:04 02.01.2006"))
	msgFormatted = fmt.Sprintf("I <b>deleted</b> the reminder \"%s\" (<i>at %s</i>) for you.", reminder.Message, reminder.RemindTime.Format("15:04 02.01.2006"))
	_, err = s.messenger.SendFormattedMessage(msg, msgFormatted, channel.ChannelIdentifier)
	return err
}
