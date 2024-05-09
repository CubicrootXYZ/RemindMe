package reply_test

import (
	"errors"
	"testing"
	"time"

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
		t.Run(msg, func(_ *testing.T) {
			event := tests.TestEvent(
				tests.MessageWithBody(
					msg,
					msg,
				))

			// Expectations
			db.EXPECT().DeleteEvent(tests.NewEventMatcher(tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()).Event)).
				Return(nil)

			msngr.EXPECT().SendResponse(gomock.Any()).Return(&messenger.MessageResponse{
				ExternalIdentifier: "abcde",
			}, nil)
			matrixDB.EXPECT().NewMessage(gomock.Any()).Return(nil, nil)

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
			action.HandleEvent(event, tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()))
			// Give async message handling some time.
			time.Sleep(time.Millisecond * 10)
		})
	}
}

func TestDeleteEventAction_HandleEventWithFailingListMessages(t *testing.T) {
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
		t.Run(msg, func(_ *testing.T) {
			event := tests.TestEvent(
				tests.MessageWithBody(
					msg,
					msg,
				))

			// Expectations
			db.EXPECT().DeleteEvent(tests.NewEventMatcher(tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()).Event)).
				Return(nil)

			msngr.EXPECT().SendResponse(gomock.Any()).Return(&messenger.MessageResponse{
				ExternalIdentifier: "abcde",
			}, nil)
			matrixDB.EXPECT().NewMessage(gomock.Any()).Return(nil, nil)

			matrixDB.EXPECT().ListMessages(matrixdb.ListMessageOpts{
				RoomID:  &tests.TestEvent().Room.ID,
				EventID: tests.TestMessage().EventID,
			}).Return([]matrixdb.MatrixMessage{}, errors.New("test"))

			// Execute
			action.HandleEvent(event, tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()))
			// Give async message handling some time.
			time.Sleep(time.Millisecond * 10)
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
		tests.MessageWithBody(
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
		tests.MessageWithBody(
			"delete",
			"delete",
		))

	// Expect
	db.EXPECT().DeleteEvent(tests.NewEventMatcher(tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()).Event)).
		Return(errors.New("test"))

	// Execute
	action.HandleEvent(event, tests.TestMessage(tests.WithFromTestEvent(), tests.WithTestEvent()))
}
