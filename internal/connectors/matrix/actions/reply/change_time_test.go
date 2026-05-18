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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangeTimeAction(t *testing.T) {
	action := &reply.ChangeTimeAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestChangeTimeAction_Selector(t *testing.T) {
	action := &reply.ChangeTimeAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()
	for _, example := range examples {
		assert.True(t, r.MatchString(example))
	}
}

func TestChangeTimeAction_HandleEvent(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &reply.ChangeTimeAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	msgs := []string{
		"my test reminder at monday 1:11",
		"my +#ä§$&7(&$==§ é reminder in 100 years",
	}

	for _, msg := range msgs {
		t.Run(msg, func(_ *testing.T) {
			event := tests.TestEvent(
				tests.MessageWithBody(
					msg,
					msg,
				))

			// Expectations
			tests.ExpectNewMessageFromEvent(matrixDB, event, matrixdb.MessageTypeChangeEvent, tests.MsgWithDBEventID(1))

			e := tests.TestMessage().Event
			e.Active = true
			db.EXPECT().UpdateEvent(mock.MatchedBy(tests.NewEventMatcher(e))).
				Return(nil, nil)

			msngr.EXPECT().SendResponse(mock.Anything).Return(&messenger.MessageResponse{
				ExternalIdentifier: "abcde",
			}, nil)
			matrixDB.EXPECT().NewMessage(mock.Anything).Return(nil, nil)

			// Execute
			action.HandleEvent(event, tests.TestMessage())
			// Wait for async message processing.
			time.Sleep(time.Millisecond * 5)
		})
	}
}

func TestChangeTimeAction_HandleEventWithUpdateError(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &reply.ChangeTimeAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	event := tests.TestEvent(
		tests.MessageWithBody(
			"tomorrow",
			"tomorrow",
		))

	// Expectations
	tests.ExpectNewMessageFromEvent(matrixDB, event, matrixdb.MessageTypeChangeEvent, tests.MsgWithDBEventID(1))

	e := tests.TestMessage().Event
	e.Active = true
	db.EXPECT().UpdateEvent(mock.MatchedBy(tests.NewEventMatcher(e))).
		Return(nil, errors.New("test"))

	// Execute
	action.HandleEvent(event, tests.TestMessage())
}

func TestChangeTimeAction_HandleEventWithNewMessageError(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &reply.ChangeTimeAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	event := tests.TestEvent(
		tests.MessageWithBody(
			"tomorrow",
			"tomorrow",
		))

	// Expectations
	matrixDB.EXPECT().NewMessage(mock.Anything).Return(nil, errors.New("test"))

	// Execute
	action.HandleEvent(event, tests.TestMessage())
}
