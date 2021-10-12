package matrixsyncer

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix/event"
)

func (s *Syncer) getActionDelete() *types.Action {
	action := &types.Action{
		Name:     "Delete all data from current user",
		Examples: []string{"delete all my data from remindme", "remove my data at remindme"},
		Regex:    "(?i)((^delete|^remove)(| all|)( my| every| the) data( from| at) remindme$)",
		Action:   s.actionDelete,
	}
	return action
}

// actionList performs the action "list" that writes all pending reminders to the given channel
func (s *Syncer) actionDelete(evt *event.Event, channel *database.Channel) error {
	msg := "Removed all your channels and data. If you have channels open I invited you into please ask the administrator to remove you from the configuration file."

	err := s.daemon.Database.DeleteChannelsFromUser(channel.UserIdentifier)
	if err != nil {
		log.Error(err.Error())
		msg = "Ups, that went wrong. Please let your administrator know that I messed that up."
	}

	_, err = s.messenger.SendFormattedMessage(msg, msg, channel, database.MessageTypeDoNotSave, 0)
	return err
}
