package message_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRegenICalTokenAction(t *testing.T) {
	action := &message.RegenICalTokenAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestRegenICalTokenAction_Selector(t *testing.T) {
	action := &message.RegenICalTokenAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()
	for _, example := range examples {
		assert.True(t, r.MatchString(example))
	}
}

func TestRegenICalTokenAction_HandleEvent(t *testing.T) { //nolint: dupl
	user := "@user:example.com"

	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.RegenICalTokenAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
		&matrix.BridgeServices{
			ICal: icalBridge,
		},
	)

	icalBridge.EXPECT().GetOutput(uint(3), true).Return(
		&icaldb.IcalOutput{
			Token: "12345",
		},
		"https://example.com/ical/1?token=abcde",
		nil,
	)
	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "evt1",
		UserID:        &user,
		RoomID:        0,
		Body:          "message",
		BodyFormatted: "<b>message</b>",
		SendAt:        time.UnixMilli(92848488),
		Type:          matrixdb.MessageTypeIcalRegenToken,
		Incoming:      true,
	})

	msngr.EXPECT().SendResponse(messenger.PlainTextResponse(
		"Your new secret calendar URL is: https://example.com/ical/1?token=abcde",
		"evt1",
		"message",
		"@user:example.com",
		"!room123",
	)).Return(&messenger.MessageResponse{
		ExternalIdentifier: "resp1",
		Timestamp:          time.UnixMilli(92848490),
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "resp1",
		UserID:        &user,
		RoomID:        0,
		Body:          "Your new secret calendar URL is: https://example.com/ical/1?token=abcde",
		BodyFormatted: "Your new secret calendar URL is: https://example.com/ical/1?token=abcde",
		SendAt:        time.UnixMilli(92848490),
		Type:          matrixdb.MessageTypeIcalRegenToken,
		Incoming:      false,
	})

	action.HandleEvent(tests.TestEvent(
		tests.WithOutput(database.Output{
			Model: gorm.Model{
				ID: 2,
			},
			OutputType: ical.OutputType,
			OutputID:   3,
		}),
	))
}

func TestRegenICalTokenAction_HandleEventWithNoOutput(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.RegenICalTokenAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
		&matrix.BridgeServices{
			ICal: icalBridge,
		},
	)

	icalBridge.EXPECT().GetOutput(uint(3), true).Return(
		nil,
		"",
		ical.ErrNotFound,
	)

	msngr.EXPECT().SendResponseAsync(messenger.PlainTextResponse(
		"It looks like iCal output is not set up for this channel. Set it up first.",
		"evt1",
		"message",
		"@user:example.com",
		"!room123",
	)).Return(nil)

	action.HandleEvent(tests.TestEvent(
		tests.WithOutput(database.Output{
			Model: gorm.Model{
				ID: 2,
			},
			OutputType: ical.OutputType,
			OutputID:   3,
		}),
	))
}
