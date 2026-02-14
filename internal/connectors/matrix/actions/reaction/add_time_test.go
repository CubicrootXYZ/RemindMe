package reaction_test

import (
	"errors"
	"log/slog"
	"slices"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/reaction"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddTimeEventAction(t *testing.T) {
	action := &reaction.AddTimeAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestAddTimeEventAction_Selector(t *testing.T) {
	action := &reaction.AddTimeAction{}

	examples := []string{}

	_, _, examplesFromDocu := action.GetDocu()
	examples = append(examples, examplesFromDocu...)

	reactions := action.Selector()

	for _, example := range examples {
		matches := slices.Contains(reactions, example)

		assert.Truef(t, matches, "%s is not part of reactions", example)
	}
}

func TestAddTimeAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &reaction.AddTimeAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	t.Run("add 1 hour", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)

		// Expectations
		db.EXPECT().UpdateEvent(gomock.Any()).Return(nil, nil)

		msngr.EXPECT().SendResponseAsync(gomock.Any()).Return(nil)

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("1️⃣"),
		), msg)
	})

	t.Run("add 1 week", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)

		// Expectations
		db.EXPECT().UpdateEvent(gomock.Any()).Return(nil, nil)

		msngr.EXPECT().SendResponseAsync(gomock.Any()).Return(nil)

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("⏩"),
		), msg)
	})

	t.Run("sending response fails", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)

		// Expectations
		db.EXPECT().UpdateEvent(gomock.Any()).Return(nil, nil)

		msngr.EXPECT().SendResponseAsync(gomock.Any()).Return(errors.New("test"))

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("2️⃣"),
		), msg)
	})

	t.Run("update fails", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)

		// Expectations
		db.EXPECT().UpdateEvent(gomock.Any()).Return(nil, errors.New("test"))

		msngr.EXPECT().SendMessageAsync(messenger.PlainTextMessage(
			"Whoopsie, can not update the event as requested.",
			"!room123",
		)).Return(nil)

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("3️⃣"),
		), msg)
	})

	t.Run("message has no event", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithoutEvent(),
		)

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("3️⃣"),
		), msg)
	})
}
