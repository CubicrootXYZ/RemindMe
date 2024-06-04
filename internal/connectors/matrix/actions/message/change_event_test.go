package message_test

import (
	"errors"
	"testing"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestChangeEventAction(t *testing.T) {
	action := &message.ChangeEventAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestChangeEventAction_Selector(t *testing.T) {
	action := &message.ChangeEventAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()
	for _, example := range examples {
		assert.True(t, r.MatchString(example))
	}
}

func TestChangeEventAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &message.ChangeEventAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	evt := database.Event{
		Model: gorm.Model{
			ID: 45,
		},
	}

	db.EXPECT().ListEvents(&database.ListEventsOpts{
		IDs:       []uint{1},
		ChannelID: &tests.TestEvent().Channel.ID,
	}).Return([]database.Event{evt}, nil)

	db.EXPECT().UpdateEvent(&eventMatcher{
		evt: &evt,
	}).Return(&evt, nil)

	msngr.EXPECT().SendMessage(messenger.HTMLMessage(
		`I rescheduled your reminder
> 
to `+today9PM(),
		`I rescheduled your reminder<br><blockquote></blockquote><br>to `+today9PM(),
		"!room123",
	)).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:     "id1",
		UserID: toP("@user:example.com"),
		Body: `I rescheduled your reminder
> 
to ` + today9PM(),
		BodyFormatted: `I rescheduled your reminder<br><blockquote></blockquote><br>to ` + today9PM(),
		Type:          matrixdb.MessageTypeChangeEvent,
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("change reminder 1 to 9 pm", "change reminder 1 to  9 pm")))

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}

func TestChangeEventAction_HandleEventWithUpdateError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &message.ChangeEventAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	evt := database.Event{
		Model: gorm.Model{
			ID: 45,
		},
	}

	db.EXPECT().ListEvents(&database.ListEventsOpts{
		IDs:       []uint{1},
		ChannelID: &tests.TestEvent().Channel.ID,
	}).Return([]database.Event{evt}, nil)

	db.EXPECT().UpdateEvent(&eventMatcher{
		evt: &evt,
	}).Return(nil, errors.New("test"))

	msngr.EXPECT().SendResponse(messenger.PlainTextResponse(
		"Whups, this did not work, sorry.",
		"evt1",
		"change reminder 1 to 9 pm",
		"@user:example.com",
		"!room123",
	)).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "id1",
		UserID:        toP("@user:example.com"),
		Body:          `Whups, this did not work, sorry.`,
		BodyFormatted: `Whups, this did not work, sorry.`,
		Type:          matrixdb.MessageTypeChangeEventError,
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("change reminder 1 to 9 pm", "change reminder 1 to  9 pm")))

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}

func TestChangeEventAction_HandleEventWithNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &message.ChangeEventAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	db.EXPECT().ListEvents(&database.ListEventsOpts{
		IDs:       []uint{1},
		ChannelID: &tests.TestEvent().Channel.ID,
	}).Return([]database.Event{}, nil)

	msngr.EXPECT().SendResponse(messenger.PlainTextResponse(
		"This reminder is not in my database.",
		"evt1",
		"change reminder 1 to 9 pm",
		"@user:example.com",
		"!room123",
	)).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "id1",
		UserID:        toP("@user:example.com"),
		Body:          `This reminder is not in my database.`,
		BodyFormatted: `This reminder is not in my database.`,
		Type:          matrixdb.MessageTypeChangeEventError,
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("change reminder 1 to 9 pm", "change reminder 1 to  9 pm")))

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}

func TestChangeEventAction_HandleEventWithMissingID(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &message.ChangeEventAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	msngr.EXPECT().SendResponse(messenger.PlainTextResponse(
		"Ups, seems like there is a reminder ID missing in your message.",
		"evt1",
		"change reminder  to abcde",
		"@user:example.com",
		"!room123",
	)).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "id1",
		UserID:        toP("@user:example.com"),
		Body:          `Ups, seems like there is a reminder ID missing in your message.`,
		BodyFormatted: `Ups, seems like there is a reminder ID missing in your message.`,
		Type:          matrixdb.MessageTypeChangeEventError,
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.MessageWithBody("change reminder  to abcde", "change reminder  to abcde")))

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}

func today9PM() string {
	now := time.Now().UTC()
	if now.Hour() >= 21 {
		now = now.Add(time.Hour * 4)
	}

	return time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, time.UTC).Format(format.DateTimeFormatDefault)
}
