package matrix

import (
	"errors"
	"testing"
	"time"

	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestService_SendDailyReminder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	fx.matrixDB.EXPECT().GetRoomByID(uint(78)).Return(
		&matrixdb.MatrixRoom{
			RoomID: "!1234",
			Model: gorm.Model{
				ID: 12,
			},
		},
		nil,
	)

	fx.messenger.EXPECT().SendMessage(messenger.HTMLMessage(
		`Your Events for Today

➡️ TEST EVENT
at 11:45 12.11.2014 (UTC) (ID: 56) 
`,
		`<h2>Your Events for Today</h2>
➡️ <b>test event</b><br>at 11:45 12.11.2014 (UTC) (ID: 56) <br>`,
		"!1234",
	)).Return(
		&messenger.MessageResponse{
			ExternalIdentifier: "abcde",
		},
		nil,
	)

	fx.matrixDB.EXPECT().NewMessage(gomock.Any()).Return(nil, nil)

	err := service.SendDailyReminder(
		&daemon.DailyReminder{
			Events: []daemon.Event{
				{
					ID:        56,
					Message:   "test event",
					EventTime: refTime(),
				},
			},
		},
		&daemon.Output{
			OutputType: "matrix",
			OutputID:   78,
		},
	)
	require.NoError(t, err)
}

func TestService_SendDailyReminderWithNewMessageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	fx.matrixDB.EXPECT().GetRoomByID(uint(78)).Return(
		&matrixdb.MatrixRoom{
			RoomID: "!1234",
			Model: gorm.Model{
				ID: 12,
			},
		},
		nil,
	)

	fx.messenger.EXPECT().SendMessage(messenger.HTMLMessage(
		`Your Events for Today

➡️ TEST EVENT
at 11:45 12.11.2014 (UTC) (ID: 56) 
`,
		`<h2>Your Events for Today</h2>
➡️ <b>test event</b><br>at 11:45 12.11.2014 (UTC) (ID: 56) <br>`,
		"!1234",
	)).Return(
		&messenger.MessageResponse{
			ExternalIdentifier: "abcde",
		},
		nil,
	)

	fx.matrixDB.EXPECT().NewMessage(gomock.Any()).Return(nil, errors.New("test"))

	err := service.SendDailyReminder(
		&daemon.DailyReminder{
			Events: []daemon.Event{
				{
					ID:        56,
					Message:   "test event",
					EventTime: refTime(),
				},
			},
		},
		&daemon.Output{
			OutputType: "matrix",
			OutputID:   78,
		},
	)
	require.NoError(t, err)
}

func TestService_SendDailyReminderWithSendMessageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	expectedErr := errors.New("test error")

	fx.matrixDB.EXPECT().GetRoomByID(uint(78)).Return(
		&matrixdb.MatrixRoom{
			RoomID: "!1234",
			Model: gorm.Model{
				ID: 12,
			},
		},
		nil,
	)

	fx.messenger.EXPECT().SendMessage(messenger.HTMLMessage(
		`Your Events for Today

➡️ TEST EVENT
at 11:45 12.11.2014 (UTC) (ID: 56) 
`,
		`<h2>Your Events for Today</h2>
➡️ <b>test event</b><br>at 11:45 12.11.2014 (UTC) (ID: 56) <br>`,
		"!1234",
	)).Return(
		nil, expectedErr,
	)

	err := service.SendDailyReminder(
		&daemon.DailyReminder{
			Events: []daemon.Event{
				{
					ID:        56,
					Message:   "test event",
					EventTime: refTime(),
				},
			},
		},
		&daemon.Output{
			OutputType: "matrix",
			OutputID:   78,
		},
	)
	require.ErrorIs(t, err, expectedErr)
}

func TestService_SendDailyReminderWithGetRoomError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, fx := testService(ctrl)

	expectedErr := errors.New("test error")

	fx.matrixDB.EXPECT().GetRoomByID(uint(78)).Return(
		nil, expectedErr,
	)

	err := service.SendDailyReminder(
		&daemon.DailyReminder{
			Events: []daemon.Event{
				{
					ID:        56,
					Message:   "test event",
					EventTime: refTime(),
				},
			},
		},
		&daemon.Output{
			OutputType: "matrix",
			OutputID:   78,
		},
	)
	require.ErrorIs(t, err, expectedErr)
}

func refTime() time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	str1 := "2014-11-12T11:45:26.371Z"
	refTime, err := time.Parse(layout, str1)
	if err != nil {
		panic(err)
	}
	return refTime
}
