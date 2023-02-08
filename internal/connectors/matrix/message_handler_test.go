package matrix

import (
	"errors"
	"regexp"
	"testing"

	"github.com/CubicrootXYZ/gologger"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type fixture struct {
	matrixDB             *matrixdb.MockService
	db                   *database.MockService
	defaultMessageAction *MockMessageAction
	messageAction        *MockMessageAction
	defaultReplyAction   *MockReplyAction
	replyAction          *MockReplyAction
}

func testService(ctrl *gomock.Controller) (service, *fixture) {
	fx := fixture{
		matrixDB:             matrixdb.NewMockService(ctrl),
		db:                   database.NewMockService(ctrl),
		defaultMessageAction: NewMockMessageAction(ctrl),
		messageAction:        NewMockMessageAction(ctrl),
		defaultReplyAction:   NewMockReplyAction(ctrl),
		replyAction:          NewMockReplyAction(ctrl),
	}

	s := service{
		config: &Config{
			DefaultMessageAction: fx.defaultMessageAction,
			MessageActions:       []MessageAction{fx.messageAction},
			DefaultReplyAction:   fx.defaultReplyAction,
			ReplyActions:         []ReplyAction{fx.replyAction},
		},
		database:       fx.db,
		matrixDatabase: fx.matrixDB,
		logger:         gologger.New(gologger.LogLevelDebug, 0),
	}

	return s, &fx
}

func TestService_MessageEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MessageEventContent{
				Body: "msg",
			},
		},
	}

	fx.matrixDB.EXPECT().GetRoomByID("abc").Return(
		&matrixdb.MatrixRoom{
			RoomID: "abc",
			Users: []matrixdb.MatrixUser{
				{
					ID: "@user:example.com",
				},
			},
		}, nil,
	)
	fx.matrixDB.EXPECT().GetMessageByID("123").Return(nil, errors.New("test"))
	fx.messageAction.EXPECT().Selector().Return(regexp.MustCompile("^$"))
	fx.defaultMessageAction.EXPECT().HandleEvent(
		&MessageEvent{
			Event:       &evt,
			Content:     evt.Content.Parsed.(*event.MessageEventContent),
			IsEncrypted: false,
		},
	)

	service.MessageEventHandler(mautrix.EventSourceTimeline, &evt)
}

func TestService_MessageEventHandlerWithMatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MessageEventContent{
				Body: "msg",
			},
		},
	}

	fx.matrixDB.EXPECT().GetRoomByID("abc").Return(
		&matrixdb.MatrixRoom{
			RoomID: "abc",
			Users: []matrixdb.MatrixUser{
				{
					ID: "@user:example.com",
				},
			},
		}, nil,
	)
	fx.matrixDB.EXPECT().GetMessageByID("123").Return(nil, errors.New("test"))
	fx.messageAction.EXPECT().Selector().Return(regexp.MustCompile("^msg$"))
	fx.messageAction.EXPECT().Name().Return("message action")
	fx.messageAction.EXPECT().HandleEvent(
		&MessageEvent{
			Event:       &evt,
			Content:     evt.Content.Parsed.(*event.MessageEventContent),
			IsEncrypted: false,
		},
	)

	service.MessageEventHandler(mautrix.EventSourceTimeline, &evt)
}

func TestService_MessageEventHandlerWithAlreadyKnown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MessageEventContent{
				Body: "msg",
			},
		},
	}

	fx.matrixDB.EXPECT().GetRoomByID("abc").Return(
		&matrixdb.MatrixRoom{
			RoomID: "abc",
			Users: []matrixdb.MatrixUser{
				{
					ID: "@user:example.com",
				},
			},
		}, nil,
	)
	fx.matrixDB.EXPECT().GetMessageByID("123").Return(&matrixdb.MatrixMessage{}, nil)

	service.MessageEventHandler(mautrix.EventSourceTimeline, &evt)
}

func TestService_MessageEventHandlerWithUserNotInRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MessageEventContent{
				Body: "msg",
			},
		},
	}

	fx.matrixDB.EXPECT().GetRoomByID("abc").Return(
		&matrixdb.MatrixRoom{
			RoomID: "abc",
			Users: []matrixdb.MatrixUser{
				{
					ID: "@user2:example.com",
				},
			},
		}, nil,
	)

	service.MessageEventHandler(mautrix.EventSourceTimeline, &evt)
}

func TestService_MessageEventHandlerWithGetRoomError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MessageEventContent{
				Body: "msg",
			},
		},
	}

	fx.matrixDB.EXPECT().GetRoomByID("abc").Return(nil, errors.New("test"))

	service.MessageEventHandler(mautrix.EventSourceTimeline, &evt)
}

func TestService_MessageEventHandlerWithDefaultReply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MessageEventContent{
				Body: "msg",
				RelatesTo: &event.RelatesTo{
					InReplyTo: &event.InReplyTo{
						EventID: id.EventID("456"),
					},
				},
			},
		},
	}

	fx.matrixDB.EXPECT().GetRoomByID("abc").Return(
		&matrixdb.MatrixRoom{
			RoomID: "abc",
			Users: []matrixdb.MatrixUser{
				{
					ID: "@user:example.com",
				},
			},
		}, nil,
	)
	fx.matrixDB.EXPECT().GetMessageByID("123").Return(nil, errors.New("test"))
	fx.matrixDB.EXPECT().GetMessageByID("456").Return(&matrixdb.MatrixMessage{}, nil)
	fx.replyAction.EXPECT().Selector().Return(regexp.MustCompile("^$"))
	fx.defaultReplyAction.EXPECT().HandleEvent(
		&MessageEvent{
			Event:       &evt,
			Content:     evt.Content.Parsed.(*event.MessageEventContent),
			IsEncrypted: false,
		},
		&matrixdb.MatrixMessage{},
	)

	service.MessageEventHandler(mautrix.EventSourceTimeline, &evt)
}

func TestService_MessageEventHandlerWithReply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MessageEventContent{
				Body: "msg",
				RelatesTo: &event.RelatesTo{
					InReplyTo: &event.InReplyTo{
						EventID: id.EventID("456"),
					},
				},
			},
		},
	}

	fx.matrixDB.EXPECT().GetRoomByID("abc").Return(
		&matrixdb.MatrixRoom{
			RoomID: "abc",
			Users: []matrixdb.MatrixUser{
				{
					ID: "@user:example.com",
				},
			},
		}, nil,
	)
	fx.matrixDB.EXPECT().GetMessageByID("123").Return(nil, errors.New("test"))
	fx.matrixDB.EXPECT().GetMessageByID("456").Return(&matrixdb.MatrixMessage{}, nil)
	fx.replyAction.EXPECT().Selector().Return(regexp.MustCompile("^msg$"))
	fx.replyAction.EXPECT().Name().Return("reply action")
	fx.replyAction.EXPECT().HandleEvent(
		&MessageEvent{
			Event:       &evt,
			Content:     evt.Content.Parsed.(*event.MessageEventContent),
			IsEncrypted: false,
		},
		&matrixdb.MatrixMessage{},
	)

	service.MessageEventHandler(mautrix.EventSourceTimeline, &evt)
}
