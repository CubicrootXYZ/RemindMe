package message_test

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	_ "time/tzdata"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaultReminderTimeAction(t *testing.T) {
	action := &message.SetDefaultReminderTimeAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestSetDefaultReminderTimeAction_Selector(t *testing.T) {
	action := &message.SetDefaultReminderTimeAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()
	for _, example := range examples {
		assert.Truef(t, r.MatchString(example), "not matching: %s", example)
	}
}

func TestSetDefaultReminderTimeAction_HandleEvent(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)
	icalBridge := ical.NewMockService(t)

	action := &message.SetDefaultReminderTimeAction{}
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

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "evt1",
		UserID:        new("@user:example.com"),
		Body:          "default reminder at 10am",
		BodyFormatted: "default reminder at 10am",
		Type:          matrixdb.MessageTypeSetDefaultReminderTime,
		Incoming:      true,
		SendAt:        time.UnixMilli(tests.TestEvent().Event.Timestamp),
	},
	).Return(nil, nil)

	channel := tests.TestEvent().Channel
	channel.DefaultReminderTime = new(uint(600))
	db.EXPECT().UpdateChannel(channel).Return(nil, nil)

	msngr.EXPECT().SendResponse(&messenger.Response{
		Message:                   "I will use 10:00 as default for your reminders.",
		MessageFormatted:          "I will use 10:00 as default for your reminders.",
		RespondToMessage:          "default reminder at 10am",
		RespondToMessageFormatted: "default reminder at 10am",
		RespondToUserID:           "@user:example.com",
		RespondToEventID:          "evt1",
		ChannelExternalIdentifier: "!room123",
	}).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "id1",
		UserID:        new("@user:example.com"),
		Body:          "I will use 10:00 as default for your reminders.",
		BodyFormatted: "I will use 10:00 as default for your reminders.",
		Type:          matrixdb.MessageTypeSetDefaultReminderTime,
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("default reminder at 10am", "default reminder at 10am")))

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}

func TestSetDefaultReminderTimeAction_HandleEvent_WithTimezone(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)
	icalBridge := ical.NewMockService(t)

	action := &message.SetDefaultReminderTimeAction{}
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

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "evt1",
		UserID:        new("@user:example.com"),
		Body:          "default reminder at 10am",
		BodyFormatted: "default reminder at 10am",
		Type:          matrixdb.MessageTypeSetDefaultReminderTime,
		Incoming:      true,
		SendAt:        time.UnixMilli(tests.TestEvent().Event.Timestamp),
	},
	).Return(nil, nil)

	channel := tests.TestEvent().Channel
	channel.DefaultReminderTime = new(uint(600))
	db.EXPECT().UpdateChannel(channel).Return(nil, nil)

	msngr.EXPECT().SendResponse(&messenger.Response{
		Message:                   "I will use 10:00 as default for your reminders.",
		MessageFormatted:          "I will use 10:00 as default for your reminders.",
		RespondToMessage:          "default reminder at 10am",
		RespondToMessageFormatted: "default reminder at 10am",
		RespondToUserID:           "@user:example.com",
		RespondToEventID:          "evt1",
		ChannelExternalIdentifier: "!room123",
	}).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "id1",
		UserID:        new("@user:example.com"),
		Body:          "I will use 10:00 as default for your reminders.",
		BodyFormatted: "I will use 10:00 as default for your reminders.",
		Type:          matrixdb.MessageTypeSetDefaultReminderTime,
	},
	).Return(nil, nil)

	event := tests.TestEvent(tests.MessageWithBody("default reminder at 10am", "default reminder at 10am"))
	event.Room.TimeZone = "Etc/GMT-2"

	action.HandleEvent(event)

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}

func TestSetDefaultReminderTimeAction_HandleEventWithUpdateChannelError(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)
	icalBridge := ical.NewMockService(t)

	action := &message.SetDefaultReminderTimeAction{}
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

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "evt1",
		UserID:        new("@user:example.com"),
		Body:          "default reminder at 10am",
		BodyFormatted: "default reminder at 10am",
		Type:          matrixdb.MessageTypeSetDefaultReminderTime,
		Incoming:      true,
		SendAt:        time.UnixMilli(tests.TestEvent().Event.Timestamp),
	},
	).Return(nil, nil)

	channel := tests.TestEvent().Channel
	channel.DefaultReminderTime = new(uint(600))
	db.EXPECT().UpdateChannel(channel).Return(nil, errors.New("test"))

	msngr.EXPECT().SendMessage(&messenger.Message{
		Body:                      "Whups, could not save that change. Sorry, try again later.",
		BodyHTML:                  "Whups, could not save that change. Sorry, try again later.",
		ChannelExternalIdentifier: "!room123",
	}).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "id1",
		UserID:        new("@user:example.com"),
		Body:          "Whups, could not save that change. Sorry, try again later.",
		BodyFormatted: "Whups, could not save that change. Sorry, try again later.",
		Type:          matrixdb.MessageTypeSetDefaultReminderTimeError,
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("default reminder at 10am", "default reminder at 10am")))

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}

func TestSetDefaultReminderTimeAction_HandleEventWithNewMessageError(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)
	icalBridge := ical.NewMockService(t)

	action := &message.SetDefaultReminderTimeAction{}
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

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "evt1",
		UserID:        new("@user:example.com"),
		Body:          "default reminder at 10am",
		BodyFormatted: "default reminder at 10am",
		Type:          matrixdb.MessageTypeSetDefaultReminderTime,
		Incoming:      true,
		SendAt:        time.UnixMilli(tests.TestEvent().Event.Timestamp),
	},
	).Return(nil, errors.New("test"))

	channel := tests.TestEvent().Channel
	channel.DefaultReminderTime = new(uint(600))
	db.EXPECT().UpdateChannel(channel).Return(nil, nil)

	msngr.EXPECT().SendResponse(&messenger.Response{
		Message:                   "I will use 10:00 as default for your reminders.",
		MessageFormatted:          "I will use 10:00 as default for your reminders.",
		RespondToMessage:          "default reminder at 10am",
		RespondToMessageFormatted: "default reminder at 10am",
		RespondToUserID:           "@user:example.com",
		RespondToEventID:          "evt1",
		ChannelExternalIdentifier: "!room123",
	}).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "id1",
		UserID:        new("@user:example.com"),
		Body:          "I will use 10:00 as default for your reminders.",
		BodyFormatted: "I will use 10:00 as default for your reminders.",
		Type:          matrixdb.MessageTypeSetDefaultReminderTime,
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("default reminder at 10am", "default reminder at 10am")))

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}
