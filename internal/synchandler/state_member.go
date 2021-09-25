package synchandler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// StateMemberHandler handles state_member events
type StateMemberHandler struct {
	database     types.Database
	messenger    types.Messenger
	matrixClient *mautrix.Client
	botInfo      *types.BotInfo
}

// NewStateMemberHandler returns a new StateMemberHandler
func NewStateMemberHandler(database types.Database, messenger types.Messenger, matrixClient *mautrix.Client, botInfo *types.BotInfo) *StateMemberHandler {
	return &StateMemberHandler{
		database:     database,
		messenger:    messenger,
		matrixClient: matrixClient,
		botInfo:      botInfo,
	}
}

// NewEvent takes a new matrix event and handles it
func (s *StateMemberHandler) NewEvent(source mautrix.EventSource, evt *event.Event) {
	content, ok := evt.Content.Parsed.(*event.MemberEventContent)
	if !ok {
		log.Warn("Event is not a member event. Can not handle it.")
		return
	}

	if evt.Timestamp/1000 < time.Now().Unix()-60 {
		return
	}

	// TODO remove
	v, _ := json.Marshal(evt)
	log.Info(fmt.Sprintf("HERE: %s", v))

	switch content.Membership {
	case event.MembershipInvite, event.MembershipJoin:
		err := s.handleInvite(evt, content)
		if err != nil {
			log.Error("Failed to handle membership invite with: " + err.Error())
		}
		return
	case event.MembershipBan, event.MembershipLeave:
		err := s.handleLeave(evt, content)
		if err != nil {
			log.Error("Failed to handle membership leave with: " + err.Error())
		}
		return
	}

	log.Info(fmt.Sprintf("No handling of this event as Membership %s is unknown.", content.Membership))
}

func (s *StateMemberHandler) handleInvite(evt *event.Event, content *event.MemberEventContent) error {
	// Ignore messages from the bot itself or if invites are not allowed
	if !s.botInfo.AllowInvites || evt.Sender.String() == s.botInfo.BotName {
		return nil
	}

	_, err := s.matrixClient.JoinRoom(evt.RoomID.String(), "", nil)
	if err != nil {
		log.Error(fmt.Sprintf("Failed joining channel %s with: %s", evt.RoomID.String(), err.Error()))
		return err
	}

	channel, err := s.database.GetChannelByUserAndChannelIdentifier(evt.Sender.String(), evt.RoomID.String())
	if err == nil && channel != nil {
		// We already know this user
		return nil
	}

	channel, err = s.database.AddChannel(evt.Sender.String(), evt.RoomID.String(), roles.RoleUser)
	if err != nil {
		return err
	}

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

	message, messageFormatted := msg.Build()

	_, err = s.messenger.SendFormattedMessage(message, messageFormatted, channel, database.MessageTypeWelcome, 0)

	return err
}

func (s *StateMemberHandler) handleLeave(evt *event.Event, content *event.MemberEventContent) error {
	channels, err := s.database.GetChannelsByChannelIdentifier(evt.RoomID.String())
	if err != nil {
		return err
	}

	for _, channel := range channels {
		err := s.database.DeleteChannel(&channel)
		if err != nil {
			log.Error("Failed to delete channel with: " + err.Error())
		}
	}

	return nil
}
