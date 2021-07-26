package matrixsyncer

import (
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

// ActionList performs the action "list" that writes all pending reminders to the given channel
func (s *Syncer) ActionList(channel *database.Channel) error {
	reminders, err := s.daemon.Database.GetPendingReminders(channel)
	if err != nil {
		return err
	}

	msg := strings.Builder{}
	msgFormatted := strings.Builder{}

	msg.WriteString("You asked for your open reminders, here they are: \n\n")
	msgFormatted.WriteString("You asked for your open reminders, here they are: <br><br>")

	for _, reminder := range reminders {
		msg.WriteString("== " + reminder.Message + " ==\n at " + reminder.RemindTime.Format("15:04 02.01.2006") + "\n\n")
		msgFormatted.WriteString("<b>" + reminder.Message + "</b><br> <i>at " + reminder.RemindTime.Format("15:04 02.01.2006") + "</i><br><br>")
	}

	_, err = s.messenger.SendFormattedMessage(msg.String(), msgFormatted.String(), channel.ChannelIdentifier)
	return err
}
