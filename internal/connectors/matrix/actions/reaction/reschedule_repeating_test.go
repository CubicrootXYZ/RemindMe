package reaction_test

import (
	"errors"
	"log/slog"
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
		matches := false
		for _, reaction := range reactions {
			if example == reaction {
				matches = true
				break
			}
		}
		assert.Truef(t, matches, "%s is not part of reactions", example)
	}
}

func TestRescheduleRepeatingAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

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
		db.EXPECT().NewEvent(gomock.Any()).Return(nil, nil)

		msngr.EXPECT().DeleteMessageAsync(gomock.Any()).Return(nil)

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
		db.EXPECT().NewEvent(gomock.Any()).Return(nil, errors.New("test"))

		msngr.EXPECT().SendMessageAsync(messenger.PlainTextMessage(
			"Whoopsie, can not update the event as requested.",
			"!room123",
		)).Return(nil)

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
		db.EXPECT().NewEvent(gomock.Any()).Return(nil, nil)

		msngr.EXPECT().DeleteMessageAsync(gomock.Any()).Return(errors.New("test"))

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("✅"),
		), msg)
	})
}
