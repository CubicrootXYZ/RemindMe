package reply_test

import (
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
	assert.NotNil(t, r)
}

func TestChangeTimeAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &reply.ChangeTimeAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
	)

	msgs := []string{
		"my test reminder at monday 1:11",
		"my +#ä§$&7(&$==§ é reminder in 100 years",
	}

	for _, msg := range msgs {
		t.Run(msg, func(t *testing.T) {
			event := tests.TestEvent(
				tests.WithBody(
					msg,
					msg,
				))

			// Expectations
			tests.ExpectNewMessageFromEvent(matrixDB, event, matrixdb.MessageTypeChangeEvent)

			db.EXPECT().UpdateEvent(tests.NewEventMatcher(tests.TestMessage().Event)).
				Return(nil, nil)

			msngr.EXPECT().SendResponseAsync(gomock.Any()).Return(nil)

			// Execute
			action.HandleEvent(event, tests.TestMessage())
		})
	}
}

// TODO more tests