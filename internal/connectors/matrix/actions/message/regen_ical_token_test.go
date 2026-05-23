package message_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
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

	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)
	icalBridge := ical.NewMockService(t)

	action := &message.RegenICalTokenAction{}
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
	}).Return(nil, nil)

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
	}).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(
		tests.MessageWithOutput(database.Output{
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
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)
	icalBridge := ical.NewMockService(t)

	action := &message.RegenICalTokenAction{}
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
		tests.MessageWithOutput(database.Output{
			Model: gorm.Model{
				ID: 2,
			},
			OutputType: ical.OutputType,
			OutputID:   3,
		}),
	))
}
