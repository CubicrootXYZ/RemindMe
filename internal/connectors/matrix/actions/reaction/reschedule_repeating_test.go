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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRescheduleRepeatAction(t *testing.T) {
	action := &reaction.RescheduleRepeatingAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestRescheduleRepeatingAction_Selector(t *testing.T) {
	action := &reaction.RescheduleRepeatingAction{}

	examples := []string{}

	_, _, examplesFromDocu := action.GetDocu()
	examples = append(examples, examplesFromDocu...)

	reactions := action.Selector()

	for _, example := range examples {
		matches := slices.Contains(reactions, example)

		assert.Truef(t, matches, "%s is not part of reactions", example)
	}
}

func TestRescheduleRepeatingAction_HandleEvent(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &reaction.RescheduleRepeatingAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	t.Run("success case", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)

		// Expectations
		db.EXPECT().NewEvent(mock.Anything).Return(nil, nil)

		msngr.EXPECT().DeleteMessageAsync(mock.Anything).Return(nil)

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("✅"),
		), msg)
	})

	t.Run("new event fails", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)

		// Expectations
		db.EXPECT().NewEvent(mock.Anything).Return(nil, errors.New("test"))

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("✅"),
		), msg)
	})

	t.Run("delete message fails", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)

		// Expectations
		db.EXPECT().NewEvent(mock.Anything).Return(nil, nil)

		msngr.EXPECT().DeleteMessageAsync(mock.Anything).Return(errors.New("test"))

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("✅"),
		), msg)
	})
}
