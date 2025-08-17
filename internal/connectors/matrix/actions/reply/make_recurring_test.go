package reply_test

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/reply"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMakeRecurringAction(t *testing.T) {
	action := &reply.MakeRecurringAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestMakeRecurringAction_Selector(t *testing.T) {
	action := &reply.MakeRecurringAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()
	for _, example := range examples {
		assert.True(t, r.MatchString(example))
	}
}

func TestMakeRecurringAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &reply.MakeRecurringAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	msg := "every 2 hours"
	event := tests.TestEvent(
		tests.MessageWithBody(
			msg,
			msg,
		))

	// Expectations
	tests.ExpectNewMessageFromEvent(matrixDB, event, matrixdb.MessageTypeChangeEvent, tests.MsgWithDBEventID(1))

	db.EXPECT().UpdateEvent(tests.NewEventMatcher(tests.TestMessage(
		tests.WithFromTestEvent(),
		tests.WithTestEvent(),
		tests.WithRecurringEvent(time.Hour*2),
	).Event)).
		Return(&database.Event{
			RepeatUntil: tests.ToP(time.Now()),
		}, nil)

	msngr.EXPECT().SendResponse(gomock.Any()).Return(&messenger.MessageResponse{}, nil)
	matrixDB.EXPECT().NewMessage(gomock.Any()).Return(nil, nil)

	// Execute
	action.HandleEvent(event, tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()))
	time.Sleep(time.Millisecond * 10)
}

func TestMakeRecurringAction_HandleEventWithUpdateError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &reply.MakeRecurringAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	msg := "every 2 hours"
	event := tests.TestEvent(
		tests.MessageWithBody(
			msg,
			msg,
		))

	// Expectations
	tests.ExpectNewMessageFromEvent(matrixDB, event, matrixdb.MessageTypeChangeEvent, tests.MsgWithDBEventID(1))

	db.EXPECT().UpdateEvent(tests.NewEventMatcher(tests.TestMessage(
		tests.WithFromTestEvent(),
		tests.WithTestEvent(),
		tests.WithRecurringEvent(time.Hour*2),
	).Event)).
		Return(nil, errors.New("test"))

	msngr.EXPECT().SendResponseAsync(gomock.Any()).Return(nil)

	// Execute
	action.HandleEvent(event, tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()))
	time.Sleep(time.Millisecond * 10)
}

func TestMakeRecurringAction_HandleEventWithDurationError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &reply.MakeRecurringAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	msg := "every 1 second"
	event := tests.TestEvent(
		tests.MessageWithBody(
			msg,
			msg,
		))

	// Expectations
	msngr.EXPECT().SendResponseAsync(gomock.Any()).Return(nil)

	// Execute
	action.HandleEvent(event, tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()))
	time.Sleep(time.Millisecond * 10)
}
