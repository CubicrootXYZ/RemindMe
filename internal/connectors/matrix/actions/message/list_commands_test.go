package message_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListCommandsAction(t *testing.T) {
	action := &message.ListCommandsAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestListCommandsAction_Selector(t *testing.T) {
	action := &message.ListCommandsAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()
	for _, example := range examples {
		assert.True(t, r.MatchString(example))
	}
}

func TestListCommandsAction_HandleEvent(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &message.ListCommandsAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	// Expectations
	matrixDB.EXPECT().NewMessage(mock.Anything).Times(3).Return(nil, nil)

	msngr.EXPECT().SendMessage(mock.Anything).Times(3).Return(&messenger.MessageResponse{
		ExternalIdentifier: "ext1",
	}, nil)

	// Execute
	action.HandleEvent(tests.TestEvent(
		tests.MessageWithBody(
			"list my commands",
			"list my commands",
		),
	))
	time.Sleep(time.Millisecond * 50) // wait for goroutine to finish
}
