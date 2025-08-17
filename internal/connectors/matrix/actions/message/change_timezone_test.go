package message_test

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestChangeTimezoneAction(t *testing.T) {
	action := &message.ChangeTimezoneAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestChangeTimezoneAction_Selector(t *testing.T) {
	action := &message.ChangeTimezoneAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()
	for _, example := range examples {
		assert.True(t, r.MatchString(example))
	}
}

func TestChangeTimezoneAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.ChangeTimezoneAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		&matrix.BridgeServices{
			ICal: icalBridge,
		},
	)

	messages := map[string]string{
		"Set timezone Europe/Berlin":       "Europe/Berlin",
		"set timezone Europe/Berlin":       "Europe/Berlin",
		"set timezone UTC":                 "UTC",
		"SeT tImezone Australia/Melbourne": "Australia/Melbourne",
	}

	for msg, tz := range messages {
		t.Run(msg, func(_ *testing.T) {
			matrixDB.EXPECT().UpdateRoom(&matrixdb.MatrixRoom{
				RoomID:   "!room123",
				Users:    []matrixdb.MatrixUser{},
				TimeZone: tz,
			}).Return(nil, nil)

			msngr.EXPECT().SendResponse(messenger.PlainTextResponse(
				"Changed this channels timezone from UTC to "+tz+" 🛫 🛬",
				"evt1",
				msg,
				"@user:example.com",
				"!room123",
			)).Return(&messenger.MessageResponse{
				ExternalIdentifier: "id1",
			}, nil)

			matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
				ID:            "id1",
				UserID:        toP("@user:example.com"),
				Body:          `Changed this channels timezone from UTC to ` + tz + ` 🛫 🛬`,
				BodyFormatted: `Changed this channels timezone from UTC to ` + tz + ` 🛫 🛬`,
				Type:          matrixdb.MessageTypeTimezoneChange,
			},
			).Return(nil, nil)

			action.HandleEvent(tests.TestEvent(tests.MessageWithBody(msg, msg)))
			// Wait for async message sending.
			time.Sleep(time.Millisecond * 10)
		})
	}
}

func TestChangeTimezoneAction_HandleEventWithInvalidTimezone(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.ChangeTimezoneAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		&matrix.BridgeServices{
			ICal: icalBridge,
		},
	)

	msngr.EXPECT().SendResponse(messenger.PlainTextResponse(
		"Sorry, but I do not know what timezone this is.",
		"evt1",
		"set timezone abc",
		"@user:example.com",
		"!room123",
	)).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "id1",
		UserID:        toP("@user:example.com"),
		Body:          `Sorry, but I do not know what timezone this is.`,
		BodyFormatted: `Sorry, but I do not know what timezone this is.`,
		Type:          matrixdb.MessageTypeTimezoneChange,
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("set timezone abc", "set timezone abc")))
}

func TestChangeTimezoneAction_HandleEventWithUpdateError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.ChangeTimezoneAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		&matrix.BridgeServices{
			ICal: icalBridge,
		},
	)

	matrixDB.EXPECT().UpdateRoom(&matrixdb.MatrixRoom{
		RoomID:   "!room123",
		Users:    []matrixdb.MatrixUser{},
		TimeZone: "Europe/Berlin",
	}).Return(nil, errors.New("test"))

	msngr.EXPECT().SendResponseAsync(messenger.PlainTextResponse(
		"Ups, that did not work 😨",
		"evt1",
		"set timezone Europe/Berlin",
		"@user:example.com",
		"!room123",
	)).Return(nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("set timezone Europe/Berlin", "set timezone Europe/Berlin")))
}
