package matrix

import (
	"errors"
	"testing"
	"time"

	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	gomock "github.com/golang/mock/gomock"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func TestService_EventStateHandlerWithInviteAndGetRoomError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipJoin,
			},
		},
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomCount().Return(int64(0), nil)
	fx.matrixDB.EXPECT().GetUserByID("@user:example.com").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(nil, errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithInviteWhitelistAndGetRoomError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)
	service.config.AllowInvites = false
	service.config.UserWhitelist = []string{"@user:example.com"}

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipJoin,
			},
		},
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetUserByID("@user:example.com").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(nil, errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithInviteAndGetUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipJoin,
			},
		},
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomCount().Return(int64(0), nil)
	fx.matrixDB.EXPECT().GetUserByID("@user:example.com").Return(nil, errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithInviteAndGetRoomCountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipJoin,
			},
		},
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomCount().Return(int64(0), errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}
func TestService_EventStateHandlerWithInviteAndRoomLimitExceeded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipJoin,
			},
		},
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomCount().Return(int64(100), nil)

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithInviteAndInvitesDisallowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)
	service.config.AllowInvites = false

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipJoin,
			},
		},
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithEventKnown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipJoin,
			},
		},
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(&matrixdb.MatrixEvent{}, nil)

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithGetEventError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipJoin,
			},
		},
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithEventFromBot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _ := testService(ctrl)

	evt := event.Event{
		Sender:    "@bot:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipJoin,
			},
		},
	}

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithEventToOld(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _ := testService(ctrl)
	service.lastMessageFrom = time.Now()

	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 1,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipJoin,
			},
		},
	}

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithLeaveAndDeleteChannelError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	stateKey := "asdfsdg"
	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipLeave,
			},
		},
		StateKey: &stateKey,
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(&matrixdb.MatrixRoom{}, nil)
	fx.matrixDB.EXPECT().DeleteAllEventsFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteAllMessagesFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().RemoveDanglingUsers().Return(int64(0), nil)
	fx.db.EXPECT().GetOutputByType(uint(0), "matrix").Return(&database.Output{ChannelID: 123}, nil)
	fx.db.EXPECT().RemoveOutputFromChannel(uint(123), uint(0)).Return(nil)
	fx.db.EXPECT().GetInputByType(uint(0), "matrix").Return(&database.Input{ChannelID: 123}, nil)
	fx.db.EXPECT().RemoveInputFromChannel(uint(123), uint(0)).Return(nil)
	fx.db.EXPECT().GetChannelByID(uint(123)).Return(&database.Channel{}, nil)
	fx.db.EXPECT().DeleteChannel(uint(0)).Return(errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithLeaveAndGetChannelError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	stateKey := "asdfsdg"
	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipLeave,
			},
		},
		StateKey: &stateKey,
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(&matrixdb.MatrixRoom{}, nil)
	fx.matrixDB.EXPECT().DeleteAllEventsFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteAllMessagesFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().RemoveDanglingUsers().Return(int64(0), nil)
	fx.db.EXPECT().GetOutputByType(uint(0), "matrix").Return(&database.Output{ChannelID: 123}, nil)
	fx.db.EXPECT().RemoveOutputFromChannel(uint(123), uint(0)).Return(nil)
	fx.db.EXPECT().GetInputByType(uint(0), "matrix").Return(&database.Input{ChannelID: 123}, nil)
	fx.db.EXPECT().RemoveInputFromChannel(uint(123), uint(0)).Return(nil)
	fx.db.EXPECT().GetChannelByID(uint(123)).Return(nil, errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithRemoveInputError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	stateKey := "asdfsdg"
	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipLeave,
			},
		},
		StateKey: &stateKey,
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(&matrixdb.MatrixRoom{}, nil)
	fx.matrixDB.EXPECT().DeleteAllEventsFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteAllMessagesFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().RemoveDanglingUsers().Return(int64(0), nil)
	fx.db.EXPECT().GetOutputByType(uint(0), "matrix").Return(&database.Output{ChannelID: 123}, nil)
	fx.db.EXPECT().RemoveOutputFromChannel(uint(123), uint(0)).Return(nil)
	fx.db.EXPECT().GetInputByType(uint(0), "matrix").Return(&database.Input{ChannelID: 123}, nil)
	fx.db.EXPECT().RemoveInputFromChannel(uint(123), uint(0)).Return(errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithGetInputError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	stateKey := "asdfsdg"
	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipLeave,
			},
		},
		StateKey: &stateKey,
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(&matrixdb.MatrixRoom{}, nil)
	fx.matrixDB.EXPECT().DeleteAllEventsFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteAllMessagesFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().RemoveDanglingUsers().Return(int64(0), nil)
	fx.db.EXPECT().GetOutputByType(uint(0), "matrix").Return(&database.Output{ChannelID: 123}, nil)
	fx.db.EXPECT().RemoveOutputFromChannel(uint(123), uint(0)).Return(nil)
	fx.db.EXPECT().GetInputByType(uint(0), "matrix").Return(nil, errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithRemoveOutputError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	stateKey := "asdfsdg"
	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipLeave,
			},
		},
		StateKey: &stateKey,
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(&matrixdb.MatrixRoom{}, nil)
	fx.matrixDB.EXPECT().DeleteAllEventsFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteAllMessagesFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().RemoveDanglingUsers().Return(int64(0), nil)
	fx.db.EXPECT().GetOutputByType(uint(0), "matrix").Return(&database.Output{ChannelID: 123}, nil)
	fx.db.EXPECT().RemoveOutputFromChannel(uint(123), uint(0)).Return(errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithGetOutputError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	stateKey := "asdfsdg"
	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipLeave,
			},
		},
		StateKey: &stateKey,
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(&matrixdb.MatrixRoom{}, nil)
	fx.matrixDB.EXPECT().DeleteAllEventsFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteAllMessagesFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().RemoveDanglingUsers().Return(int64(0), nil)
	fx.db.EXPECT().GetOutputByType(uint(0), "matrix").Return(nil, errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithDeleteRoomError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	stateKey := "asdfsdg"
	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipLeave,
			},
		},
		StateKey: &stateKey,
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(&matrixdb.MatrixRoom{}, nil)
	fx.matrixDB.EXPECT().DeleteAllEventsFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteAllMessagesFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteRoom(uint(0)).Return(errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithDeleteMessagesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	stateKey := "asdfsdg"
	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipLeave,
			},
		},
		StateKey: &stateKey,
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(&matrixdb.MatrixRoom{}, nil)
	fx.matrixDB.EXPECT().DeleteAllEventsFromRoom(uint(0)).Return(nil)
	fx.matrixDB.EXPECT().DeleteAllMessagesFromRoom(uint(0)).Return(errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithDeleteEventsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	stateKey := "asdfsdg"
	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipLeave,
			},
		},
		StateKey: &stateKey,
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(&matrixdb.MatrixRoom{}, nil)
	fx.matrixDB.EXPECT().DeleteAllEventsFromRoom(uint(0)).Return(errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}

func TestService_EventStateHandlerWithGetRoomError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	stateKey := "asdfsdg"
	evt := event.Event{
		Sender:    "@user:example.com",
		RoomID:    id.RoomID("abc"),
		Timestamp: 5000,
		ID:        id.EventID("123"),
		Content: event.Content{
			Parsed: &event.MemberEventContent{
				Membership: event.MembershipLeave,
			},
		},
		StateKey: &stateKey,
	}

	fx.matrixDB.EXPECT().GetEventByID("123").Return(nil, matrixdb.ErrNotFound)
	fx.matrixDB.EXPECT().GetRoomByRoomID("abc").Return(nil, errors.New("test"))

	service.EventStateHandler(mautrix.EventSourceAccountData, &evt)
}
