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

func TestDeleteEventAction(t *testing.T) {
	action := &reaction.DeleteEventAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestDeleteEventAction_Selector(t *testing.T) {
	action := &reaction.DeleteEventAction{}

	examples := []string{}

	_, _, examplesFromDocu := action.GetDocu()
	examples = append(examples, examplesFromDocu...)

	reactions := action.Selector()

	for _, example := range examples {
		matches := slices.Contains(reactions, example)

		assert.Truef(t, matches, "%s is not part of reactions", example)
	}
}

func TestDeleteAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &reaction.DeleteEventAction{}
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
		evt := *msg.Event

		// Expectations
		evt.Active = false
		db.EXPECT().UpdateEvent(&evt).Return(nil, nil)

		msngr.EXPECT().SendResponseAsync(messenger.PlainTextResponse(
			"Deleted event \"\"",
			"msg1",
			"",
			"@user:example.com",
			"!room123",
		)).Return(nil)

		matrixDB.EXPECT().ListMessages(matrixdb.ListMessageOpts{
			RoomID:  &tests.TestEvent().Room.ID,
			EventID: tests.TestMessage().EventID,
		}).Return([]matrixdb.MatrixMessage{
			{
				ID: "123456",
			},
			{
				ID: "1234567",
			},
		}, nil)

		msngr.EXPECT().DeleteMessageAsync(&messenger.Delete{
			ExternalIdentifier:        "123456",
			ChannelExternalIdentifier: tests.TestEvent().Room.RoomID,
		}).Return(nil)

		msngr.EXPECT().DeleteMessageAsync(&messenger.Delete{
			ExternalIdentifier:        "1234567",
			ChannelExternalIdentifier: tests.TestEvent().Room.RoomID,
		}).Return(errors.New("test"))

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("❌"),
		), msg)
	})

	t.Run("sending response fails", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)
		evt := *msg.Event

		// Expectations
		evt.Active = false
		db.EXPECT().UpdateEvent(&evt).Return(nil, nil)

		msngr.EXPECT().SendResponseAsync(messenger.PlainTextResponse(
			"Deleted event \"\"",
			"msg1",
			"",
			"@user:example.com",
			"!room123",
		)).Return(errors.New("test"))

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("❌"),
		), msg)
	})

	t.Run("update fails", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithTestEvent(),
		)
		evt := *msg.Event

		// Expectations
		evt.Active = false
		db.EXPECT().UpdateEvent(&evt).Return(nil, errors.New("test"))

		msngr.EXPECT().SendMessageAsync(messenger.PlainTextMessage(
			"Whoopsie, can not delete the event as requested.",
			"!room123",
		)).Return(nil)

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("❌"),
		), msg)
	})

	t.Run("missing event in message", func(_ *testing.T) {
		msg := tests.TestMessage(
			tests.WithoutEvent(),
		)

		// Execute
		action.HandleEvent(tests.TestReactionEvent(
			tests.ReactionWithKey("❌"),
		), msg)
	})
}
