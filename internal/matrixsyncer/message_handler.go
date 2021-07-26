package matrixsyncer

import (
	"fmt"
	"regexp"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"gorm.io/gorm"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// handles new messages
func (s *Syncer) handleMessages(source mautrix.EventSource, evt *event.Event) {
	log.Debug(fmt.Sprintf("New message: / Sender: %s / Room: / %s / Time: %d", evt.Sender, evt.RoomID, evt.Timestamp))

	// Do not answer our own and old messages
	if evt.Sender == id.UserID(s.botName) || evt.Timestamp/1000 <= time.Now().Unix()-60 {
		return
	}

	channel, err := s.daemon.Database.GetChannelByUserAndChannelIdentifier(evt.Sender.String(), evt.RoomID.String())

	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		log.Warn("Event is not a message event. Can not handle it")
		return
	}

	// Unknown channel
	if err == gorm.ErrRecordNotFound || channel == nil {
		channel2, _ := s.daemon.Database.GetChannelByUserIdentifier(evt.Sender.String())
		// But we know the user
		if channel2 != nil {
			log.Info("User messaged us in a Channel we do not know")
			_, err := s.messenger.SendReplyToEvent("Hey, this is not our usual messaging channel ;)", evt, evt.RoomID.String())
			if err != nil {
				log.Warn(err.Error())
			}
		} else {
			log.Info("We do not know that user.")
		}
		return
	}

	// TODO handle replies

	if s.checkActions(evt, channel, content) {
		return
	}

	// Nothing left so it must be a reminder
	reminder, err := s.parseRemind(evt, channel)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed parsing the Reminder with: %s", err.Error()))
		return
	}

	msg := fmt.Sprintf("Successfully added new reminder (ID: %d) for %s", reminder.ID, reminder.RemindTime.Format("15:04 02.01.2006"))

	log.Info(msg)
	_, err = s.messenger.SendReplyToEvent(msg, evt, evt.RoomID.String())
	if err != nil {
		log.Warn("Was not able to send success message to user")
	}
}

// checkActions checks if a message matches any special actions and performs them.
func (s *Syncer) checkActions(evt *event.Event, channel *database.Channel, content *event.MessageEventContent) (matched bool) {
	// List action
	if matched, err := regexp.Match("(?i)((^list|^show)(| all| the)(| reminders| my reminders)(| please)$|^reminders$|^reminder$)", []byte(content.Body)); matched && err == nil {
		s.ActionList(channel)
		return true
	}

	return false
}
