package matrixsyncer

import (
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getActionList() *Action {
	action := &Action{
		Name:     "List all reminders",
		Examples: []string{"list", "list reminders", "show", "show reminders", "list my reminders", "reminders"},
		Regex:    "(?i)((^list|^show)(| all| the)(| reminders| my reminders)(| please)$|^reminders$|^reminder$)",
		Action:   s.actionList,
	}
	return action
}

// actionList performs the action "list" that writes all pending reminders to the given channel
func (s *Syncer) actionList(evt *event.Event, channel *database.Channel) error {
	reminders, err := s.daemon.Database.GetPendingReminders(channel)
	if err != nil {
		return err
	}

	msg := strings.Builder{}
	msgFormatted := strings.Builder{}

	msg.WriteString("= Open Reminders = \nYou asked for your open reminders, here they are: \n\n")
	msgFormatted.WriteString("<h3>Open Reminders</h3><br>You asked for your open reminders, here they are: <br><br>")

	for _, reminder := range reminders {
		msg.WriteString("== " + reminder.Message + " ==\n at " + reminder.RemindTime.Format("15:04 02.01.2006") + "\n\n")
		msgFormatted.WriteString("<b>" + reminder.Message + "</b><br> <i>at " + reminder.RemindTime.Format("15:04 02.01.2006") + "</i><br><br>")
	}

	_, err = s.messenger.SendFormattedMessage(msg.String(), msgFormatted.String(), channel.ChannelIdentifier)
	return err
}
