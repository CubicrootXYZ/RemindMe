package message_test

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

func TestAddUserAction(t *testing.T) {
	action := &message.AddUserAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestAddUserAction_Selector(t *testing.T) {
	action := &message.AddUserAction{}

	shouldMatch := []string{
		"  add   user    @mybuddy:matrix.org",
		"add user " + format.GetMatrixLinkForUser("@user"),
	}

	_, _, examples := action.GetDocu()
	shouldMatch = append(shouldMatch, examples...)

	r := action.Selector()
	for _, msg := range shouldMatch {
		assert.Truef(t, r.MatchString(msg), "'%s' should match but did not", msg)
	}
}

func TestAddUserAction_HandleEvent(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &message.AddUserAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	t.Run("formatted username", func(_ *testing.T) { //nolint:dupl
		// Expectations
		client.EXPECT().JoinedMembers(id.RoomID("!room123")).Return(
			&mautrix.RespJoinedMembers{
				Joined: map[id.UserID]mautrix.JoinedMember{
					id.UserID("@user:example.org"): {},
				},
			},
			nil,
		)
		matrixDB.EXPECT().AddUserToRoom("@user:example.org", tests.TestEvent().Room).
			Return(nil, nil)
		msngr.EXPECT().SendResponse(mock.Anything).Return(&messenger.MessageResponse{
			ExternalIdentifier: "resp1",
			Timestamp:          time.UnixMilli(92848490),
		}, nil)
		matrixDB.EXPECT().NewMessage(mock.Anything).Return(nil, nil)
		matrixDB.EXPECT().NewMessage(mock.Anything).Return(nil, nil)

		// Execute
		action.HandleEvent(tests.TestEvent(
			tests.MessageWithBody(
				"add user @user:example.org",
				`add user <a href="https://matrix.to/#/@user:example.org" class="linkified" rel="noreferrer noopener">@user:example.org</a>`,
			),
		))
	})

	t.Run("plain text username", func(_ *testing.T) { //nolint:dupl
		// Expectations
		client.EXPECT().JoinedMembers(id.RoomID("!room123")).Return(
			&mautrix.RespJoinedMembers{
				Joined: map[id.UserID]mautrix.JoinedMember{
					id.UserID("@user:example.org"): {},
				},
			},
			nil,
		)
		matrixDB.EXPECT().AddUserToRoom("@user:example.org", tests.TestEvent().Room).
			Return(nil, nil)
		msngr.EXPECT().SendResponse(mock.Anything).Return(&messenger.MessageResponse{
			ExternalIdentifier: "resp1",
			Timestamp:          time.UnixMilli(92848490),
		}, nil)
		matrixDB.EXPECT().NewMessage(mock.Anything).Return(nil, nil)
		matrixDB.EXPECT().NewMessage(mock.Anything).Return(nil, nil)

		// Execute
		action.HandleEvent(tests.TestEvent(
			tests.MessageWithBody(
				"add user @user:example.org",
				"",
			),
		))
	})
}

func TestAddUserAction_HandleEventWithResponseFailed(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &message.AddUserAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	// Expectations
	client.EXPECT().JoinedMembers(id.RoomID("!room123")).Return(
		&mautrix.RespJoinedMembers{
			Joined: map[id.UserID]mautrix.JoinedMember{
				id.UserID("@user:example.org"): {},
			},
		},
		nil,
	)
	matrixDB.EXPECT().AddUserToRoom("@user:example.org", tests.TestEvent().Room).
		Return(nil, nil)
	matrixDB.EXPECT().NewMessage(mock.Anything).Return(nil, nil)
	msngr.EXPECT().SendResponse(mock.Anything).Return(nil, errors.New("test"))

	// Execute
	action.HandleEvent(tests.TestEvent(
		tests.MessageWithBody(
			"add user @user:example.org",
			`add user <a href="https://matrix.to/#/@user:example.org" class="linkified" rel="noreferrer noopener">@user:example.org</a>`,
		),
	))
}

func TestAddUserAction_HandleEventWithAddUserError(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &message.AddUserAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	// Expectations
	client.EXPECT().JoinedMembers(id.RoomID("!room123")).Return(
		&mautrix.RespJoinedMembers{
			Joined: map[id.UserID]mautrix.JoinedMember{
				id.UserID("@user:example.org"): {},
			},
		},
		nil,
	)
	matrixDB.EXPECT().AddUserToRoom("@user:example.org", tests.TestEvent().Room).
		Return(nil, errors.New("test"))

	// Execute
	action.HandleEvent(tests.TestEvent(
		tests.MessageWithBody(
			"add user @user:example.org",
			`add user <a href="https://matrix.to/#/@user:example.org" class="linkified" rel="noreferrer noopener">@user:example.org</a>`,
		),
	))
}

func TestAddUserAction_HandleEventWithUserAlreadyInRoom(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &message.AddUserAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	// Expectations
	client.EXPECT().JoinedMembers(id.RoomID("!room123")).Return(
		&mautrix.RespJoinedMembers{
			Joined: map[id.UserID]mautrix.JoinedMember{
				id.UserID("@user:example.org"): {},
			},
		},
		nil,
	)
	msngr.EXPECT().SendResponseAsync(mock.Anything).Return(errors.New("test"))

	// Execute
	action.HandleEvent(tests.TestEvent(
		tests.MessageWithBody(
			"add user @user:example.org",
			`add user <a href="https://matrix.to/#/@user:example.org" class="linkified" rel="noreferrer noopener">@user:example.org</a>`,
		),
		tests.MessageWithUserInRoom(
			matrixdb.MatrixUser{
				ID: "@user:example.org",
			},
		),
	))
}

func TestAddUserAction_HandleEventWithNoUsername(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &message.AddUserAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	// Expectations
	client.EXPECT().JoinedMembers(id.RoomID("!room123")).Return(
		&mautrix.RespJoinedMembers{
			Joined: map[id.UserID]mautrix.JoinedMember{
				id.UserID("@user:example.org"): {},
			},
		},
		nil,
	)
	msngr.EXPECT().SendResponseAsync(mock.Anything).Return(nil)

	// Execute
	action.HandleEvent(tests.TestEvent(
		tests.MessageWithBody(
			"",
			"",
		),
	))
}

func TestAddUserAction_HandleEventWithUserNotJoined(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &message.AddUserAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	// Expectations
	client.EXPECT().JoinedMembers(id.RoomID("!room123")).Return(
		&mautrix.RespJoinedMembers{
			Joined: map[id.UserID]mautrix.JoinedMember{},
		},
		nil,
	)
	msngr.EXPECT().SendResponseAsync(mock.Anything).Return(errors.New("test"))

	// Execute
	action.HandleEvent(tests.TestEvent(
		tests.MessageWithBody(
			"add user @user:example.org",
			`add user <a href="https://matrix.to/#/@user:example.org" class="linkified" rel="noreferrer noopener">@user:example.org</a>`,
		),
	))
}

func TestAddUserAction_HandleEventWithJoinedError(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &message.AddUserAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	// Expectations
	client.EXPECT().JoinedMembers(id.RoomID("!room123")).
		Return(nil, errors.New("test"))

	// Execute
	action.HandleEvent(tests.TestEvent(
		tests.MessageWithBody(
			"add user @user:example.org",
			`add user <a href="https://matrix.to/#/@user:example.org" class="linkified" rel="noreferrer noopener">@user:example.org</a>`,
		),
	))
}
