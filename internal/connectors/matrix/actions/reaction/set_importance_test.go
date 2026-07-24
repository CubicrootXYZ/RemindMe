package reaction_test

import (
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

func TestSetImportanceEventAction(t *testing.T) {
	action := &reaction.SetImportanceAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestSetImportanceEventAction_Selector(t *testing.T) {
	action := &reaction.SetImportanceAction{}

	examples := []string{}

	_, _, examplesFromDocu := action.GetDocu()
	examples = append(examples, examplesFromDocu...)

	reactions := action.Selector()

	for _, example := range examples {
		matches := slices.Contains(reactions, example)

		assert.Truef(t, matches, "%s is not part of reactions", example)
	}
}

func TestSetImportanceAction_HandleEvent(t *testing.T) {
	// Setup
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)
	client := mautrixcl.NewMockClient(t)
	msngr := messenger.NewMockMessenger(t)

	action := &reaction.SetImportanceAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		nil,
	)

	t.Run("set to important", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)

		// Expectations
		db.EXPECT().UpdateEvent(mock.MatchedBy(func(evt *database.Event) bool {
			return evt.Importance == database.ImportanceImportant
		})).Return(nil, nil)

		msngr.EXPECT().SendResponseAsync(mock.Anything).Return(nil)

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("⚠️"),
		), msg)
	})

	t.Run("set to default", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)

		// Expectations
		db.EXPECT().UpdateEvent(mock.MatchedBy(func(evt *database.Event) bool {
			return evt.Importance == database.ImportanceDefault
		})).Return(nil, nil)

		msngr.EXPECT().SendResponseAsync(mock.Anything).Return(nil)

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("ℹ️"),
		), msg)
	})

	t.Run("unknown key", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)

		// Expectations
		msngr.EXPECT().SendMessageAsync(mock.Anything).Return(nil)

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("❎"),
		), msg)
	})
}
