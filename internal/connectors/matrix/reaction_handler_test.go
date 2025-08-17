package matrix

import (
	"testing"

	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	gomock "github.com/golang/mock/gomock"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func TestService_ReactionEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.ReactionEventContent{
				RelatesTo: event.RelatesTo{
					Key:     "😈",
					EventID: "123",
				},
			},
		},
	}

	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(
		&matrixdb.MatrixRoom{
			RoomID: "abc",
			Users: []matrixdb.MatrixUser{
				{
					ID: "@user:example.com",
				},
			},
		}, nil,
	)
	fx.matrixDB.EXPECT().GetMessageByID("123").Return(&matrixdb.MatrixMessage{
		Room: matrixdb.MatrixRoom{
			RoomID: "abc",
		},
	}, nil)
	fx.reactionAction.EXPECT().Selector().Return([]string{"😈"})
	fx.db.EXPECT().GetInputByType(uint(0), "matrix").Return(&database.Input{}, nil)
	fx.db.EXPECT().GetChannelByID(uint(0)).Return(&database.Channel{}, nil)
	fx.reactionAction.EXPECT().HandleEvent(
		&ReactionEvent{
			Event:   &evt,
			Content: evt.Content.Parsed.(*event.ReactionEventContent),
			Room:    testRoom(),
			Channel: &database.Channel{},
			Input:   &database.Input{},
		},
		&matrixdb.MatrixMessage{
			RoomID: testRoom().ID,
			Room: matrixdb.MatrixRoom{
				RoomID: "abc",
			},
		},
	)

	service.ReactionEventHandler(mautrix.EventSourceTimeline, &evt)
}
