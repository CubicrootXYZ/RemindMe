package matrixsyncer

import (
	"regexp"
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/id"
)

func (s *Syncer) getActionAddUser() *types.Action {
	action := &types.Action{
		Name:     "Add user to interact with the bot",
		Examples: []string{"add user @bestbuddy"},
		Regex:    regexp.MustCompile("(?i)(^add user).*"),
		Action:   s.actionAddUser,
	}
	return action
}

// actionAddUser adds a user to a channel
func (s *Syncer) actionAddUser(evt *types.MessageEvent, channel *database.Channel) error {
	users, err := s.client.JoinedMembers(id.RoomID(channel.ChannelIdentifier))
	if err != nil {
		return err
	}

	msg := "Sorry, could not find that user in this channel"

	username := getUsernameFromLink(evt.Content.FormattedBody)
	exactMatch := true
	if username == "" {
		// Fall back to plain text
		username = getUsernameFromText(evt.Content.Body)
		exactMatch = false
	}

	if username == "" {
		msg := "Sorry :(, I was not able to get a user out of your message"
		err = s.messenger.SendMessageAsync(asyncmessenger.PlainTextMessage(msg, channel.ChannelIdentifier))
		return err
	}

	addUser := false
	for user := range users.Joined {
		if exactMatch {
			if user.String() == username {
				addUser = true
				break
			}
		} else {
			if "@"+username == strings.Split(user.String(), ":")[0] {
				addUser = true
				break
			}
		}
	}

	_, err = s.daemon.Database.GetChannelByUserAndChannelIdentifier(username, channel.ChannelIdentifier)
	if err == nil {
		msg = "User is already added"
		addUser = false
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	if addUser {
		_, err = s.daemon.Database.AddChannel(username, channel.ChannelIdentifier, roles.RoleUser)
		if err != nil {
			msg := "Sorry, sonething went wrong here"
			err = s.messenger.SendMessageAsync(asyncmessenger.PlainTextMessage(msg, channel.ChannelIdentifier))
			return err
		}

		form := formater.Formater{}
		form.Text("Added ")
		form.Username(username)
		form.Text(" to the channel")
		msg, msgFormatted := form.Build()
		err = s.messenger.SendMessageAsync(asyncmessenger.HTMLMessage(msg, msgFormatted, channel.ChannelIdentifier))
		return err
	}

	err = s.messenger.SendResponseAsync(asyncmessenger.PlainTextResponse(msg, string(evt.Event.ID), evt.Content.Body, evt.Event.Sender.String(), evt.Event.RoomID.String()))

	return err
}

func getUsernameFromLink(link string) string {
	r := regexp.MustCompile(`https:\/\/matrix.to\/#\/[^"'>]+`)

	url := r.Find([]byte(link))
	if url == nil {
		return ""
	}

	return strings.TrimPrefix(string(url), "https://matrix.to/#/")
}

func getUsernameFromText(text string) string {
	return strings.TrimPrefix(text, "add user ")
}
