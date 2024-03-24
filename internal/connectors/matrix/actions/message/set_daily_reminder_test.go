package message_test

import (
	"errors"
	"testing"
	"time"

	"github.com/CubicrootXYZ/gologger"
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

func TestSetDailyReminderAction(t *testing.T) {
	action := &message.SetDailyReminderAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestSetDailyReminderAction_Selector(t *testing.T) {
	action := &message.SetDailyReminderAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()
	for _, example := range examples {
		assert.True(t, r.MatchString(example))
	}
}

func TestSetDailyReminderAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.SetDailyReminderAction{}
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

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "evt1",
		UserID:        toP("@user:example.com"),
		Body:          "daily reminder at 10am",
		BodyFormatted: "daily reminder at 10am",
		Type:          matrixdb.MessageTypeSetDailyReminder,
		Incoming:      true,
		SendAt:        time.UnixMilli(tests.TestEvent().Event.Timestamp),
	},
	).Return(nil, nil)

	channel := tests.TestEvent().Channel
	channel.DailyReminder = toP(uint(600))
	db.EXPECT().UpdateChannel(channel).Return(nil, nil)

	msngr.EXPECT().SendResponse(&messenger.Response{
		Message:                   "I will send you a daily overview at 10:00. To disable the reminder message me with \"delete daily reminder\".",
		MessageFormatted:          "I will send you a daily overview at 10:00. To disable the reminder message me with \"delete daily reminder\".",
		RespondToMessage:          "daily reminder at 10am",
		RespondToMessageFormatted: "daily reminder at 10am",
		RespondToUserID:           "@user:example.com",
		RespondToEventID:          "evt1",
		ChannelExternalIdentifier: "!room123",
	}).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "id1",
		UserID:        toP("@user:example.com"),
		Body:          "I will send you a daily overview at 10:00. To disable the reminder message me with \"delete daily reminder\".",
		BodyFormatted: "I will send you a daily overview at 10:00. To disable the reminder message me with \"delete daily reminder\".",
		Type:          matrixdb.MessageTypeSetDailyReminder,
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("daily reminder at 10am", "daily reminder at 10am")))

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}

func TestSetDailyReminderAction_HandleEventWithUpdateChannelError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.SetDailyReminderAction{}
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

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "evt1",
		UserID:        toP("@user:example.com"),
		Body:          "daily reminder at 10am",
		BodyFormatted: "daily reminder at 10am",
		Type:          matrixdb.MessageTypeSetDailyReminder,
		Incoming:      true,
		SendAt:        time.UnixMilli(tests.TestEvent().Event.Timestamp),
	},
	).Return(nil, nil)

	channel := tests.TestEvent().Channel
	channel.DailyReminder = toP(uint(600))
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
		UserID:        toP("@user:example.com"),
		Body:          "Whups, could not save that change. Sorry, try again later.",
		BodyFormatted: "Whups, could not save that change. Sorry, try again later.",
		Type:          matrixdb.MessageTypeSetDailyReminderError,
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("daily reminder at 10am", "daily reminder at 10am")))

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}

func TestSetDailyReminderAction_HandleEventWithNewMessageError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.SetDailyReminderAction{}
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

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "evt1",
		UserID:        toP("@user:example.com"),
		Body:          "daily reminder at 10am",
		BodyFormatted: "daily reminder at 10am",
		Type:          matrixdb.MessageTypeSetDailyReminder,
		Incoming:      true,
		SendAt:        time.UnixMilli(tests.TestEvent().Event.Timestamp),
	},
	).Return(nil, errors.New("test"))

	channel := tests.TestEvent().Channel
	channel.DailyReminder = toP(uint(600))
	db.EXPECT().UpdateChannel(channel).Return(nil, nil)

	msngr.EXPECT().SendResponse(&messenger.Response{
		Message:                   "I will send you a daily overview at 10:00. To disable the reminder message me with \"delete daily reminder\".",
		MessageFormatted:          "I will send you a daily overview at 10:00. To disable the reminder message me with \"delete daily reminder\".",
		RespondToMessage:          "daily reminder at 10am",
		RespondToMessageFormatted: "daily reminder at 10am",
		RespondToUserID:           "@user:example.com",
		RespondToEventID:          "evt1",
		ChannelExternalIdentifier: "!room123",
	}).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "id1",
		UserID:        toP("@user:example.com"),
		Body:          "I will send you a daily overview at 10:00. To disable the reminder message me with \"delete daily reminder\".",
		BodyFormatted: "I will send you a daily overview at 10:00. To disable the reminder message me with \"delete daily reminder\".",
		Type:          matrixdb.MessageTypeSetDailyReminder,
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("daily reminder at 10am", "daily reminder at 10am")))

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}
