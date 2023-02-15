package message

import (
	"regexp"
	"strings"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix"
)

type AddUserAction struct {
	logger    gologger.Logger
	client    *mautrix.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
}

func (action *AddUserAction) Configure(logger gologger.Logger, client *mautrix.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
}

func (action *AddUserAction) Name() string {
	return "Add user"
}

func (action *AddUserAction) Selector() *regexp.Regexp {
	return regexp.MustCompile("(?i)(^add user).*")
}

func (action *AddUserAction) HandleEvent(event *matrix.MessageEvent) {
	usersInRoom, err := action.client.JoinedMembers(event.Event.RoomID)
	if err != nil {
		action.logger.Err(err)
		return
	}

	username := format.GetUsernameFromLink(event.Content.FormattedBody)
	exactMatch := true
	if username == "" {
		// Fall back to plain text
		username = strings.TrimPrefix(strings.TrimSpace(event.Content.Body), "add user ")
		exactMatch = false
	}

	// Return if username not found in message
	if username == "" {
		err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Sorry 😕, but I was not able to find a username in that message.",
			event.Event.ID.String(),
			event.Content.Body,
			event.Event.Sender.String(),
			event.Room.RoomID,
		))
		if err != nil {
			action.logger.Err(err)
		}
		return
	}

	userInRoom := false
	for user := range usersInRoom.Joined {
		if exactMatch {
			if user.String() == username {
				userInRoom = true
				break
			}
		} else {
			if "@"+username == strings.Split(user.String(), ":")[0] {
				userInRoom = true
				username = user.String()
				break
			}
		}
	}

	// Return if user is not in room
	if !userInRoom {
		err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
			"Bad news 😰, can not find that user in this room.",
			event.Event.ID.String(),
			event.Content.Body,
			event.Event.Sender.String(),
			event.Room.RoomID,
		))
		if err != nil {
			action.logger.Err(err)
		}
		return
	}

	// Return if user is already added
	for _, user := range event.Room.Users {
		if user.ID == username {
			err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
				"This user is already added.",
				event.Event.ID.String(),
				event.Content.Body,
				event.Event.Sender.String(),
				event.Room.RoomID,
			))
			if err != nil {
				action.logger.Err(err)
			}
			return
		}
	}

	// Add new user to room
	_, err = action.matrixDB.AddUserToRoom(username, event.Room)
	if err != nil {
		action.logger.Err(err)
		return
	}

	err = action.messenger.SendResponseAsync(messenger.PlainTextResponse(
		"Added that user 👏. They can now interact with me.",
		event.Event.ID.String(),
		event.Content.Body,
		event.Event.Sender.String(),
		event.Room.RoomID,
	))
	if err != nil {
		action.logger.Err(err)
	}

	// TODO add message to database
}
