package matrix

import (
	"errors"
	"time"

	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	db "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// EventStateHandler handles state events from matrix.
func (service *service) EventStateHandler(_ mautrix.EventSource, evt *event.Event) {
	logger := service.logger.With(
		"matrix.sender", evt.Sender,
		"matrix.room", evt.RoomID,
		"matrix.event_timestamp", evt.Timestamp,
	)
	logger.Debug("new state event")

	service.metricEventInCount.
		WithLabelValues("state").
		Inc()

	// Ignore old events or events from the bot itself
	if evt.Sender.String() == service.botname || evt.Timestamp/1000 <= service.lastMessageFrom.Unix() {
		return
	}

	content, ok := evt.Content.Parsed.(*event.MemberEventContent)
	if !ok {
		logger.Info("cannot handle event", "reason", "not a member event")
		return
	}

	// Check if the event is known
	_, err := service.matrixDatabase.GetEventByID(evt.ID.String())
	if err == nil {
		return
	} else if !errors.Is(err, matrixdb.ErrNotFound) {
		logger.Error("failed to get event", "error", err)
		return
	}

	switch content.Membership {
	case event.MembershipInvite, event.MembershipJoin:
		err := service.handleInvite(evt, content)
		if err != nil {
			logger.Error("failed to handle state event", "matrix.event.membership", "invite/join", "error", err)
		}
	case event.MembershipLeave, event.MembershipBan:
		err := service.handleLeave(evt)
		if err != nil {
			logger.Error("failed to handle state event", "matrix.event.membership", "leave/ban", "error", err)
		}
	default:
		logger.Info("cannot handle state event", "reason", "unknown membership type", "event.membership", content.Membership)
	}
}

func (service *service) handleInvite(evt *event.Event, content *event.MemberEventContent) error {
	declineInvites, err := service.maxUserReached()
	if err != nil {
		return err
	}

	if declineInvites && !service.userInWhitelist(evt.Sender.String()) {
		service.logger.Debug("invite ignored", "reason", "disabled/not whitelisted", "matrix.user", evt.Sender.String())
		return nil
	}

	user, err := service.matrixDatabase.GetUserByID(evt.Sender.String())
	if err != nil {
		if !errors.Is(err, matrixdb.ErrNotFound) {
			return err
		}
		user = nil
	}

	if user != nil && user.Blocked {
		service.logger.Debug("invite ignored", "reason", "blocked", "matrix.user", evt.Sender.String())
		return nil
	}

	_, err = service.matrixDatabase.GetRoomByRoomID(evt.RoomID.String())
	if err != nil {
		if !errors.Is(err, matrixdb.ErrNotFound) {
			return err
		}
	} else {
		// Room already known, ignore it
		return nil
	}

	// TODO for further testing service.client needs to be mocked.
	_, err = service.client.JoinRoom(evt.RoomID.String(), "", nil)
	if err != nil {
		service.logger.Error("failed to join channel", "error", err, "matrix.room.id", evt.RoomID.String())
		return err
	}

	if user == nil {
		user, err = service.matrixDatabase.NewUser(&matrixdb.MatrixUser{
			ID: evt.Sender.String(),
		})
		if err != nil {
			return err
		}
	}

	room, err := service.matrixDatabase.NewRoom(&matrixdb.MatrixRoom{
		RoomID: evt.RoomID.String(),
	})
	if err != nil {
		return err
	}

	room.Users = append(room.Users, *user)
	room, err = service.matrixDatabase.UpdateRoom(room)
	if err != nil {
		return err
	}

	_, err = service.matrixDatabase.NewEvent(&matrixdb.MatrixEvent{
		ID:     evt.ID.String(),
		UserID: user.ID,
		RoomID: room.ID,
		Type:   string(content.Membership),
		SendAt: time.Unix(evt.Timestamp/1000, 0),
	})
	if err != nil {
		return err
	}

	err = service.setupNewChannel(room, user)
	if err != nil {
		service.logger.Error("failed to setup new channel", "error", err)
		return err
	}

	return nil
}

func (service *service) setupNewChannel(room *matrixdb.MatrixRoom, user *matrixdb.MatrixUser) error {
	channel, err := service.database.NewChannel(&db.Channel{
		Description: "auto generated channel for matrix room " + room.RoomID,
	})
	if err != nil {
		return err
	}

	err = service.database.AddInputToChannel(
		channel.ID,
		&db.Input{
			InputType: InputType,
			InputID:   room.ID,
			Enabled:   true,
		},
	)
	if err != nil {
		return err
	}

	err = service.database.AddOutputToChannel(
		channel.ID,
		&db.Output{
			OutputType: OutputType,
			OutputID:   room.ID,
			Enabled:    true,
		},
	)
	if err != nil {
		return err
	}

	go service.sendWelcomeMessage(room, user)

	return nil
}

func (service *service) sendWelcomeMessage(room *matrixdb.MatrixRoom, user *matrixdb.MatrixUser) {
	message, messageFormatted := getWelcomeMessage(room)

	resp, err := service.messenger.SendMessage(messenger.HTMLMessage(
		message,
		messageFormatted,
		room.RoomID,
	))
	if err != nil {
		service.logger.Info("failed to send message", "error", err.Error())
		return
	}

	_, err = service.matrixDatabase.NewMessage(&matrixdb.MatrixMessage{
		ID:            resp.ExternalIdentifier,
		UserID:        &user.ID,
		RoomID:        room.ID,
		Body:          message,
		BodyFormatted: messageFormatted,
		SendAt:        resp.Timestamp,
		Type:          matrixdb.MessageTypeWelcome,
		Incoming:      false,
	})
	if err != nil {
		service.logger.Error("failed saving message to database", "error", err.Error())
	}
}

func getWelcomeMessage(room *matrixdb.MatrixRoom) (string, string) {
	msg := format.Formater{}
	msg.Title("Welcome to RemindMe")
	msg.TextLine("Hey, I am your personal reminder bot. Beep boop beep.")
	msg.Text("You want to now what I am capable of? Just text me ")
	msg.BoldLine("list all commands")
	msg.Text("Is this your current local time? ")
	msg.Italic(format.ToLocalTime(time.Now(), room.TimeZone))
	msg.NewLine()
	msg.TextLine("If not, please adjust your timezone with ")
	msg.BoldLine("set timezone Europe/Berlin")

	msg.TextLine("You can set up a daily reminder too!")

	msg.SubTitle("Attribution")
	msg.TextLine("This bot is open for everyone and build with the help of voluntary software developers.")
	msg.Text("The source code can be found at ")
	msg.Link("GitHub", "https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot")
	msg.TextLine(". Star it if you like the bot, open issues or discussions with your findings.")

	return msg.Build()
}

func (service *service) maxUserReached() (bool, error) {
	if !service.config.AllowInvites {
		return true, nil
	}

	if service.config.RoomLimit > 0 {
		roomCount, err := service.matrixDatabase.GetRoomCount()
		if err != nil {
			return true, err
		}

		if service.config.RoomLimit > uint(roomCount) {
			return false, nil
		}

		return true, nil
	}

	return false, nil
}

func (service *service) handleLeave(evt *event.Event) error {
	if evt.StateKey == nil {
		return nil
	}

	room, err := service.matrixDatabase.GetRoomByRoomID(string(evt.RoomID))
	if err != nil {
		if errors.Is(err, matrixdb.ErrNotFound) {
			return nil
		}
		return err
	}

	if time.Unix(evt.Timestamp/1000, 0).Sub(room.CreatedAt) < 0 {
		// Got invited to room after this event. Ignore this event.
		service.logger.Info("leave ignored", "reason", "got invited again", "matrix.room.id", room.RoomID)
		return nil
	}

	err = service.removeRoom(room)
	if err != nil {
		return err
	}

	err = service.removeFromChannel(room)
	if err != nil {
		service.logger.Error("failed to remove room from channel", "error", err)
		return err
	}

	// TODO mock service.client to test further.
	_, err = service.client.LeaveRoom(evt.RoomID)
	if err != nil {
		// Fire and forget, we might already be banned
		service.logger.Error("failed to leave room", "error", err)
	}

	return nil
}

func (service *service) removeRoom(room *matrixdb.MatrixRoom) error {
	err := service.matrixDatabase.DeleteAllEventsFromRoom(room.ID)
	if err != nil {
		return err
	}

	err = service.matrixDatabase.DeleteAllMessagesFromRoom(room.ID)
	if err != nil {
		return err
	}

	err = service.matrixDatabase.DeleteRoom(room.ID)
	if err != nil {
		return err
	}

	cnt, err := service.matrixDatabase.RemoveDanglingUsers()
	if err != nil {
		return err
	}

	service.logger.Info("deleted dangling matrix users", "count", cnt)

	return nil
}

func (service *service) removeFromChannel(room *matrixdb.MatrixRoom) error {
	output, err := service.database.GetOutputByType(room.ID, OutputType)
	if err == nil {
		err = service.database.RemoveOutputFromChannel(output.ChannelID, output.ID)
		if err != nil {
			return err
		}
	} else if !errors.Is(err, db.ErrNotFound) {
		return err
	}

	input, err := service.database.GetInputByType(room.ID, InputType)
	if err == nil {
		err = service.database.RemoveInputFromChannel(input.ChannelID, input.ID)
		if err != nil {
			return err
		}
	} else if !errors.Is(err, db.ErrNotFound) {
		return err
	}

	channel, err := service.database.GetChannelByID(input.ChannelID)
	if err == nil && len(channel.Inputs) == 0 && len(channel.Outputs) == 0 {
		err = service.database.DeleteChannel(channel.ID)
		if err != nil {
			return err
		}
	} else if !errors.Is(err, db.ErrNotFound) {
		return err
	}

	return nil
}

func (service *service) userInWhitelist(user string) bool {
	for _, u := range service.config.UserWhitelist {
		if u == user {
			return true
		}
	}

	return false
}
