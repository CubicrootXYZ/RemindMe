package synchandler

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"gorm.io/gorm"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// MessageHandler handles message events
type MessageHandler struct {
	database     types.Database
	messenger    types.Messenger
	botInfo      *types.BotInfo
	replyActions []*types.ReplyAction
	actions      []*types.Action
}

// NewMessageHandler returns a new MessageHandler
func NewMessageHandler(database types.Database, messenger types.Messenger, botInfo *types.BotInfo, replyActions []*types.ReplyAction, messageAction []*types.Action) *MessageHandler {
	return &MessageHandler{
		database:     database,
		messenger:    messenger,
		botInfo:      botInfo,
		replyActions: replyActions,
		actions:      messageAction,
	}
}

// NewEvent takes a new matrix event and handles it
func (s *MessageHandler) NewEvent(source mautrix.EventSource, evt *event.Event) {
	log.Debug(fmt.Sprintf("New message: / Sender: %s / Room: / %s / Time: %d", evt.Sender, evt.RoomID, evt.Timestamp))

	// Do not answer our own and old messages
	if evt.Sender == id.UserID(s.botInfo.BotName) || evt.Timestamp/1000 <= time.Now().Unix()-60 {
		return
	}
	// TODO check if the message is already known

	channel, err := s.database.GetChannelByUserAndChannelIdentifier(evt.Sender.String(), evt.RoomID.String())

	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		log.Warn("Event is not a message event. Can not handle it")
		return
	}

	// Unknown channel
	if err == gorm.ErrRecordNotFound || channel == nil {
		channel2, _ := s.database.GetChannelByUserIdentifier(evt.Sender.String())
		// But we know the user
		if channel2 != nil {
			log.Info("User messaged us in a Channel we do not know")
			_, err := s.messenger.SendReplyToEvent("Hey, this is not our usual messaging channel ;)", evt, &database.Channel{ChannelIdentifier: evt.RoomID.String()}, database.MessageTypeDoNotSave)
			if err != nil {
				log.Warn(err.Error())
			}
		} else {
			log.Info("We do not know that user.")
		}
		return
	}

	// Check if it is a reply to a message we know
	if s.checkReplyActions(evt, channel, content) {
		return
	}

	// Check if a action matches
	if s.checkActions(evt, channel, content) {
		return
	}

	// Nothing left so it must be a reminder
	_, err = s.newReminder(evt, channel)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed parsing the Reminder with: %s", err.Error()))
		return
	}
}

func (s *MessageHandler) checkReplyActions(evt *event.Event, channel *database.Channel, content *event.MessageEventContent) (matched bool) {
	if content.RelatesTo == nil || channel == nil {
		return false
	}
	if len(content.RelatesTo.EventID) < 2 {
		return false
	}

	message := strings.ToLower(formater.StripReply(content.Body))
	replyMessage, err := s.database.GetMessageByExternalID(content.RelatesTo.EventID.String())
	if err != nil || replyMessage == nil {
		log.Info("Message replies to unknown message")
		return false
	}

	// Cycle through all registered actions
	for _, action := range s.replyActions {
		log.Debug("Checking for match with " + action.Name)
		log.Debug(string(replyMessage.Type))

		for _, rtt := range action.ReplyToTypes {
			if rtt == replyMessage.Type {
				log.Debug("Regex matching: " + message)
				if matched, err := regexp.Match(action.Regex, []byte(message)); matched && err == nil {
					_ = action.Action(evt, channel, replyMessage, content)
					log.Debug("Matched")
					return true
				}
			}
		}
	}

	// Fallback change reminder date
	if replyMessage.ReminderID != nil && *replyMessage.ReminderID > 0 {
		err = s.changeReminderDate(replyMessage, channel, content, evt)
		if err != nil {
			log.Error(err.Error())
		}
		return true
	}

	return false
}

func (s *MessageHandler) changeReminderDate(replyMessage *database.Message, channel *database.Channel, content *event.MessageEventContent, evt *event.Event) error {
	remindTime, err := formater.ParseTime(content.Body, channel, false)
	if err != nil {
		log.Warn(err.Error())
		s.messenger.SendReplyToEvent("Sorry I was not able to get a time out of that message", evt, channel, database.MessageTypeReminderUpdateFail)
		return err
	}

	reminder, err := s.database.UpdateReminder(*replyMessage.ReminderID, remindTime, 0, 0)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	_, err = s.database.AddMessageFromMatrix(evt.ID.String(), evt.Timestamp, content, reminder, database.MessageTypeReminderUpdate, channel)
	if err != nil {
		log.Warn(fmt.Sprintf("Could not register reply message %s in database", evt.ID.String()))
	}

	s.messenger.SendReplyToEvent(fmt.Sprintf("I rescheduled your reminder \"%s\" to %s.", reminder.Message, formater.ToLocalTime(reminder.RemindTime, channel)), evt, channel, database.MessageTypeReminderUpdateSuccess)

	return nil
}

// checkActions checks if a message matches any special actions and performs them.
func (s *MessageHandler) checkActions(evt *event.Event, channel *database.Channel, content *event.MessageEventContent) (matched bool) {
	message := strings.ToLower(content.Body)

	// List action
	for _, action := range s.actions {
		log.Debug("Checking for match with action " + action.Name)
		if matched, err := regexp.Match(action.Regex, []byte(message)); matched && err == nil {
			_ = action.Action(evt, channel)
			log.Debug("Matched")
			return true
		}
	}

	return false
}

func (s *MessageHandler) newReminder(evt *event.Event, channel *database.Channel) (*database.Reminder, error) {
	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		return nil, errors.ErrMatrixEventWrongType
	}

	remindTime, err := formater.ParseTime(content.Body, channel, false)
	if err != nil {
		s.messenger.SendReplyToEvent("Sorry I was not able to understand the remind date and time from this message", evt, channel, database.MessageTypeReminderFail)
		return nil, err
	}

	reminder, err := s.database.AddReminder(remindTime, content.Body, true, uint64(0), channel)
	if err != nil {
		log.Warn("Error when inserting reminder: " + err.Error())
		return reminder, err
	}
	_, err = s.database.AddMessageFromMatrix(evt.ID.String(), evt.Timestamp/1000, content, reminder, database.MessageTypeReminderRequest, channel)
	if err != nil {
		log.Warn("Was not able to save a message to the database: " + err.Error())
	}

	msg := fmt.Sprintf("Successfully added new reminder (ID: %d) for %s", reminder.ID, formater.ToLocalTime(reminder.RemindTime, channel))

	log.Info(msg)
	_, err = s.messenger.SendReplyToEvent(msg, evt, channel, database.MessageTypeReminderSuccess)
	if err != nil {
		log.Warn("Was not able to send success message to user")
	}

	return reminder, err
}
