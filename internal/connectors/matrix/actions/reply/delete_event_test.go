package reply_test

import (
	"errors"
	"testing"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/reply"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteEventAction(t *testing.T) {
	action := &reply.DeleteEventAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestDeleteEventAction_Selector(t *testing.T) {
	action := &reply.DeleteEventAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()
	for _, example := range examples {
		assert.True(t, r.MatchString(example))
	}
}

func TestDeleteEventAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &reply.DeleteEventAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	msgs := []string{
		"delete",
		" remove ",
	}

	for _, msg := range msgs {
		t.Run(msg, func(t *testing.T) {
			event := tests.TestEvent(
				tests.WithBody(
					msg,
					msg,
				))

			// Expectations
			db.EXPECT().DeleteEvent(tests.NewEventMatcher(tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()).Event)).
				Return(nil)

			msngr.EXPECT().SendResponseAsync(gomock.Any()).Return(nil)

			// Execute
			action.HandleEvent(event, tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()))
		})
	}
}

func TestDeleteEventAction_HandleEventWithMissingEventID(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &reply.DeleteEventAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	event := tests.TestEvent(
		tests.WithBody(
			"delete",
			"delete",
		))

	// Execute
	action.HandleEvent(event, tests.TestMessage(tests.WithoutEvent()))
}

func TestDeleteEventAction_HandleEventWithDeleteError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &reply.DeleteEventAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	event := tests.TestEvent(
		tests.WithBody(
			"delete",
			"delete",
		))

	// Expect
	db.EXPECT().DeleteEvent(tests.NewEventMatcher(tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()).Event)).
		Return(errors.New("test"))

	// Execute
	action.HandleEvent(event, tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()))
}
