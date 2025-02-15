package api_test

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func testTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2106-01-02T15:04:05+07:00")

	return t
}

func TestAPI_ICALExportHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, icalDB, server := testServer(ctrl)

	icalDB.EXPECT().GetIcalOutputByID(uint(1)).Return(
		&icaldb.IcalOutput{
			Token: "1234",
			Model: gorm.Model{
				ID: 2,
			},
		},
		nil,
	)

	db.EXPECT().GetOutputByType(uint(2), ical.OutputType).Return(
		&database.Output{
			ChannelID: 34,
		},
		nil,
	)

	db.EXPECT().GetEventsByChannel(uint(34)).Return(
		[]database.Event{
			{
				Model: gorm.Model{
					ID: 69,
				},
				Time:    testTime(),
				Message: "my msg",
			},
		},
		nil,
	)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/ical/1?token=1234",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:34\nMETHOD:PUBLISH\nBEGIN:VEVENT\nDTSTART:21060102T080405Z\nDTEND:21060102T080905Z\nDTSTAMP:00010101T000000Z\nUID:69\nSUMMARY:my msg\nDESCRIPTION:my msg\nCLASS:PRIVATE\nEND:VEVENT\nEND:VCALENDAR\n", string(body))
}

func TestAPI_ICALExportHandlerWithEventsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, icalDB, server := testServer(ctrl)

	icalDB.EXPECT().GetIcalOutputByID(uint(1)).Return(
		&icaldb.IcalOutput{
			Token: "1234",
			Model: gorm.Model{
				ID: 2,
			},
		},
		nil,
	)

	db.EXPECT().GetOutputByType(uint(2), ical.OutputType).Return(
		&database.Output{
			ChannelID: 34,
		},
		nil,
	)

	db.EXPECT().GetEventsByChannel(uint(34)).Return(
		nil,
		errors.New("test"),
	)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/ical/1?token=1234",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAPI_ICALExportHandlerWithOutputNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, icalDB, server := testServer(ctrl)

	icalDB.EXPECT().GetIcalOutputByID(uint(1)).Return(
		&icaldb.IcalOutput{
			Token: "1234",
			Model: gorm.Model{
				ID: 2,
			},
		},
		nil,
	)

	db.EXPECT().GetOutputByType(uint(2), ical.OutputType).Return(
		nil,
		database.ErrNotFound,
	)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/ical/1?token=1234",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAPI_ICALExportHandlerWithWrongToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, icalDB, server := testServer(ctrl)

	icalDB.EXPECT().GetIcalOutputByID(uint(1)).Return(
		&icaldb.IcalOutput{
			Token: "12345",
			Model: gorm.Model{
				ID: 2,
			},
		},
		nil,
	)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/ical/1?token=1234",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAPI_ICALExportHandlerWithIcalOutputNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, icalDB, server := testServer(ctrl)

	icalDB.EXPECT().GetIcalOutputByID(uint(1)).Return(
		nil,
		icaldb.ErrNotFound,
	)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/ical/1?token=1234",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAPI_ICALExportHandlerWithNoToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, _, server := testServer(ctrl)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/ical/1",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
