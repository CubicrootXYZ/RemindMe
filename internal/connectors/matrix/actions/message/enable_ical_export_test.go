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

func TestEnableICalExportAction(t *testing.T) {
	action := &message.EnableICalExportAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestNEnableICalExportAction_Selector(t *testing.T) {
	action := &message.EnableICalExportAction{}

	r := action.Selector()
	assert.NotNil(t, r)
}

func TestEnableICalExportAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.EnableICalExportAction{}
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

	icalBridge.EXPECT().GetOutput(uint(2)).Return(
		&icaldb.IcalOutput{
			Token: "12345",
		},
		nil,
	)
	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "evt1",
		UserID:        "@user:example.com",
		RoomID:        0,
		Body:          "message",
		BodyFormatted: "<b>message</b>",
		SendAt:        time.UnixMilli(92848488),
		Type:          matrixdb.MessageTypeIcalExportEnable,
		Incoming:      true,
	})

	msngr.EXPECT().SendResponse(messenger.PlainTextResponse(
		"Your calendar is ready 🥳: 12345",
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
		UserID:        "@user:example.com",
		RoomID:        0,
		Body:          "Your calendar is ready 🥳: 12345",
		BodyFormatted: "Your calendar is ready 🥳: 12345",
		SendAt:        time.UnixMilli(92848490),
		Type:          matrixdb.MessageTypeIcalExportEnable,
		Incoming:      false,
	})

	action.HandleEvent(tests.TestEvent(
		tests.WithOutput(database.Output{
			Model: gorm.Model{
				ID: 2,
			},
			OutputType: ical.OutputType,
		}),
	))
}

func TestEnableICalExportAction_HandleEventWithNoOutput(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.EnableICalExportAction{}
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

	icalBridge.EXPECT().GetOutput(uint(2)).Return(
		nil,
		ical.ErrNotFound,
	)
	icalBridge.EXPECT().NewOutput(uint(68272)).Return(&icaldb.IcalOutput{
		Token: "12345",
	},
		nil,
	)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "evt1",
		UserID:        "@user:example.com",
		RoomID:        0,
		Body:          "message",
		BodyFormatted: "<b>message</b>",
		SendAt:        time.UnixMilli(92848488),
		Type:          matrixdb.MessageTypeIcalExportEnable,
		Incoming:      true,
	})

	msngr.EXPECT().SendResponse(messenger.PlainTextResponse(
		"Your calendar is ready 🥳: 12345",
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
		UserID:        "@user:example.com",
		RoomID:        0,
		Body:          "Your calendar is ready 🥳: 12345",
		BodyFormatted: "Your calendar is ready 🥳: 12345",
		SendAt:        time.UnixMilli(92848490),
		Type:          matrixdb.MessageTypeIcalExportEnable,
		Incoming:      false,
	})

	action.HandleEvent(tests.TestEvent(
		tests.WithOutput(database.Output{
			Model: gorm.Model{
				ID: 2,
			},
			OutputType: ical.OutputType,
		}),
	))
}
