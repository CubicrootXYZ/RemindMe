package message_test

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	dbtests "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database/tests"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListEventsAction_Meta(t *testing.T) {
	action := &message.ListEventsAction{}
	assert.Greater(t, len(action.Name()), 2)

	title, expl, examples := action.GetDocu()
	assert.Greater(t, len(title), 2)
	assert.Greater(t, len(expl), 2)
	assert.NotEmpty(t, examples)
}

func TestListEventsAction_Selector(t *testing.T) {
	action := &message.ListEventsAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()

	for _, example := range examples {
		assert.True(t, r.MatchString(example))
	}
}

func TestListEventsAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.ListEventsAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		&matrix.BridgeServices{
			ICal: icalBridge,
		},
	)

	db.EXPECT().ListEvents(&database.ListEventsOpts{
		ChannelID: new(uint(68272)),
	}).Return([]database.Event{
		dbtests.TestEvent(),
	}, nil)

	matrixDB.EXPECT().NewMessage(tests.TestMessage(
		tests.WithFromTestEvent(),
		tests.WithMessageType(matrixdb.MessageTypeEventList),
	)).Return(nil, nil)

	msngr.EXPECT().SendMessage(&messenger.Message{
		Body: `== YOUR EVENTS ==

JANUARY 2006
➡️ TEST EVENT
at 08:04 02.01.2006 (UTC) (ID: 2824) 
`,
		BodyHTML: `<h3>Your Events</h3><br><br><b>January 2006</b><br>
➡️ <b>test event</b><br>at 08:04 02.01.2006 (UTC) (ID: 2824) <br>`,
		ChannelExternalIdentifier: "!room123",
	}).Return(&messenger.MessageResponse{
		ExternalIdentifier: "!234",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:     "!234",
		UserID: new("@user:example.com"),
		Body: `== YOUR EVENTS ==

JANUARY 2006
➡️ TEST EVENT
at 08:04 02.01.2006 (UTC) (ID: 2824) 
`,
		BodyFormatted: `<h3>Your Events</h3><br><br><b>January 2006</b><br>
➡️ <b>test event</b><br>at 08:04 02.01.2006 (UTC) (ID: 2824) <br>`,
		Type:   matrixdb.MessageTypeEventList,
		RoomID: 0,
	})

	action.HandleEvent(tests.TestEvent())

	// Wait for async process.
	time.Sleep(time.Millisecond * 10)
}

func TestListEventsAction_HandleEventWithError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.ListEventsAction{}
	action.Configure(
		slog.Default(),
		client,
		msngr,
		matrixDB,
		db,
		&matrix.BridgeServices{
			ICal: icalBridge,
		},
	)

	db.EXPECT().ListEvents(&database.ListEventsOpts{
		ChannelID: new(uint(68272)),
	}).Return(nil, errors.New("test"))

	msngr.EXPECT().SendResponseAsync(&messenger.Response{
		Message:                   "There was an issue accessing the data 🤨",
		MessageFormatted:          "There was an issue accessing the data 🤨",
		RespondToMessage:          "message",
		RespondToMessageFormatted: "message",
		RespondToUserID:           "@user:example.com",
		RespondToEventID:          "evt1",
		ChannelExternalIdentifier: "!room123",
	}).Return(nil)

	action.HandleEvent(tests.TestEvent())
}
