package synchandler

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"gorm.io/gorm"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"
)

// StateMemberHandler handles state_member events
type StateMemberHandler struct {
	database     types.Database
	messenger    asyncmessenger.Messenger
	matrixClient *mautrix.Client
	botInfo      *types.BotInfo
	botSettings  *configuration.BotSettings
	olm          *crypto.OlmMachine
	started      int64
}

// NewStateMemberHandler returns a new StateMemberHandler
func NewStateMemberHandler(database types.Database, messenger asyncmessenger.Messenger, matrixClient *mautrix.Client, botInfo *types.BotInfo, botSettings *configuration.BotSettings, olm *crypto.OlmMachine) *StateMemberHandler {
	return &StateMemberHandler{
		database:     database,
		messenger:    messenger,
		matrixClient: matrixClient,
		botInfo:      botInfo,
		botSettings:  botSettings,
		olm:          olm,
		started:      time.Now().Unix(),
	}
}

// NewEvent takes a new matrix event and handles it
func (s *StateMemberHandler) NewEvent(_ mautrix.EventSource, evt *event.Event) {
	if s.olm != nil {
		s.olm.HandleMemberEvent(evt)
	}

	if evt.Timestamp/1000 < s.started {
		return
	}

	content, ok := evt.Content.Parsed.(*event.MemberEventContent)
	if !ok {
		log.Warn("Event is not a member event. Can not handle it.")
		return
	}

	// Check if the event is known
	known, err := s.database.IsEventKnown(evt.ID.String())
	if known {
		return
	}
	if err != nil {
		log.Error(err.Error())
	}

	switch content.Membership {
	case event.MembershipInvite, event.MembershipJoin:
		err := s.handleInvite(evt, content)
		if err != nil {
			log.Error("Failed to handle membership invite with: " + err.Error())
		}
	case event.MembershipLeave, event.MembershipBan:
		err := s.handleLeave(evt, content)
		if err != nil {
			log.Error("Failed to handle membership leave with: " + err.Error())
		}
	default:
		log.Info(fmt.Sprintf("No handling of this event as Membership %s is unknown.", content.Membership))
	}
}

func (s *StateMemberHandler) handleInvite(evt *event.Event, content *event.MemberEventContent) error {
	// Ignore messages from the bot itself
	if evt.Sender.String() == s.botInfo.BotName {
		return nil
	}

	declineInvites, err := s.maxUserReached()
	if err != nil {
		return err
	}
	isUserBlocked, err := s.database.IsUserBlocked(evt.Sender.String())
	if err != nil {
		return err
	}

	channels, err := s.database.GetChannelsByChannelIdentifier(evt.RoomID.String())
	if err != nil {
		return err
	}

	if len(channels) > 0 {
		// Only allow one user per channel to be auto added, others can than be added manually
		declineInvites = true
	}

	if declineInvites || isUserBlocked {
		log.Info(evt.Sender.String() + " is blocked or bot reached max users")
		return nil
	}

	_, err = s.matrixClient.JoinRoom(evt.RoomID.String(), "", nil)
	if err != nil {
		log.Error(fmt.Sprintf("Failed joining channel %s with: %s", evt.RoomID.String(), err.Error()))
		return err
	}

	channel, err := s.database.GetChannelByUserAndChannelIdentifier(evt.Sender.String(), evt.RoomID.String())
	if err == nil && channel != nil {
		// We already know this user
		return s.addMemberEventToDatabase(evt, content)
	}

	channel, err = s.database.AddChannel(evt.Sender.String(), evt.RoomID.String(), roles.RoleUser)
	if err != nil {
		return err
	}

	go func(channel *database.Channel) {
		message, messageFormatted := getWelcomeMessage()

		resp, err := s.messenger.SendMessage(asyncmessenger.HTMLMessage(
			message,
			messageFormatted,
			channel.ChannelIdentifier,
		))
		if err != nil {
			log.Info("Failed to send message: " + err.Error())
			return
		}

		_, err = s.database.AddMessage(&database.Message{
			Body:               message,
			BodyHTML:           messageFormatted,
			Type:               database.MessageTypeWelcome,
			ChannelID:          channel.ID,
			Timestamp:          resp.Timestamp,
			ExternalIdentifier: resp.ExternalIdentifier,
		})
		if err != nil {
			log.Info("Failed saving message into database: " + err.Error())
		}
	}(channel)

	err = s.addMemberEventToDatabase(evt, content)

	return err
}

func (s *StateMemberHandler) handleLeave(evt *event.Event, content *event.MemberEventContent) error {
	if evt.StateKey == nil {
		return nil
	}

	channel, err := s.database.GetChannelByUserAndChannelIdentifier(*evt.StateKey, evt.RoomID.String())
	if err != nil {
		return err
	}

	err = s.database.DeleteChannel(channel)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("Failed to delete channel with: " + err.Error())
	}

	err = s.addMemberEventToDatabase(evt, content)

	return err
}

func (s *StateMemberHandler) addMemberEventToDatabase(evt *event.Event, content *event.MemberEventContent) error {
	dbEvent := database.Event{}
	dbEvent.ExternalIdentifier = evt.ID.String()

	if content.Membership == event.MembershipInvite || content.Membership == event.MembershipJoin {
		channel, err := s.database.GetChannelByUserAndChannelIdentifier(evt.Sender.String(), evt.RoomID.String())
		if err != nil {
			return err
		}

		dbEvent.Channel = *channel
		dbEvent.ChannelID = &channel.ID
	}

	dbEvent.Timestamp = evt.Timestamp / 1000
	dbEvent.EventType = database.EventTypeMembership
	dbEvent.EventSubType = string(content.Membership)
	_, err := s.database.AddEvent(&dbEvent)

	return err
}

func getWelcomeMessage() (string, string) {
	msg := formater.Formater{}
	msg.Title("Welcome to RemindMe")
	msg.TextLine("Hey, I am your personal reminder bot. Beep boop beep.")
	msg.Text("You want to now what I am capable of? Just text me ")
	msg.BoldLine("list all commands")
	msg.TextLine("First things you should do are setting your timezone and a daily reminder.")

	msg.SubTitle("Attribution")
	msg.TextLine("This bot is open for everyone and build with the help of voluntary software developers.")
	msg.Text("The source code can be found at ")
	msg.Link("GitHub", "https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot")
	msg.TextLine(". Star it if you like the bot, open issues or discussions with your findings.")

	return msg.Build()
}

func (s *StateMemberHandler) maxUserReached() (bool, error) {
	if !s.botSettings.AllowInvites {
		return true, nil
	}

	if s.botSettings.MaxUser >= 0 {
		channelCount, err := s.database.ChannelCount()
		if err != nil {
			return true, err
		}

		if channelCount >= s.botSettings.MaxUser {
			log.Info("Reached max channels - will no longer follow new invites.")
			return true, nil
		}
	}

	return false, nil
}
