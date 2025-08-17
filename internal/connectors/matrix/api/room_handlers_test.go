package api_test

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRoom() matrixdb.MatrixRoom {
	room := matrixdb.MatrixRoom{
		RoomID: "roomid",
		Users: []matrixdb.MatrixUser{
			{
				ID: "userid",
			},
		},
	}

	room.ID = 1
	room.CreatedAt, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")

	return room
}

func TestAPI_ListInputRoomsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, matrixDB, server := testServer(ctrl)

	matrixDB.EXPECT().ListInputRoomsByChannel(uint(1)).Return(
		[]matrixdb.MatrixRoom{testRoom()},
		nil,
	)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/matrix/channels/1/inputs/rooms",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"status":"success","data":[{"id":1,"created_at":"2006-01-02T15:04:05+07:00","room_id":"roomid","users":[{"id":"userid","blocked":false}]}]}`, string(body))
}

func TestAPI_ListInputRoomsHandlerWithEncryption(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, matrixDB, server := testServer(ctrl)

	room := testRoom()
	matrixDB.EXPECT().ListInputRoomsByChannel(uint(1)).Return(
		[]matrixdb.MatrixRoom{room},
		nil,
	)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/matrix/channels/1/inputs/rooms",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"status":"success","data":[{"id":1,"created_at":"2006-01-02T15:04:05+07:00","room_id":"roomid","users":[{"id":"userid","blocked":false}]}]}`, string(body))
}

func TestAPI_ListInputRoomsHandlerWithDatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, matrixDB, server := testServer(ctrl)

	matrixDB.EXPECT().ListInputRoomsByChannel(uint(1)).Return(
		nil,
		errors.New("test"),
	)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/matrix/channels/1/inputs/rooms",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"message":"Internal Server Error","status":"error"}`, string(body))
}

func TestAPI_ListInputRoomsHandlerWithInvalidPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, _, server := testServer(ctrl)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/matrix/channels/0/inputs/rooms",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAPI_ListOutputRoomsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, matrixDB, server := testServer(ctrl)

	matrixDB.EXPECT().ListOutputRoomsByChannel(uint(1)).Return(
		[]matrixdb.MatrixRoom{testRoom()},
		nil,
	)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/matrix/channels/1/outputs/rooms",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"status":"success","data":[{"id":1,"created_at":"2006-01-02T15:04:05+07:00","room_id":"roomid","users":[{"id":"userid","blocked":false}]}]}`, string(body))
}

func TestAPI_ListOutputRoomsHandlerWithEncryption(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, matrixDB, server := testServer(ctrl)

	room := testRoom()
	matrixDB.EXPECT().ListOutputRoomsByChannel(uint(1)).Return(
		[]matrixdb.MatrixRoom{room},
		nil,
	)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/matrix/channels/1/outputs/rooms",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"status":"success","data":[{"id":1,"created_at":"2006-01-02T15:04:05+07:00","room_id":"roomid","users":[{"id":"userid","blocked":false}]}]}`, string(body))
}

func TestAPI_ListOutputRoomsHandlerWithDatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, matrixDB, server := testServer(ctrl)

	matrixDB.EXPECT().ListOutputRoomsByChannel(uint(1)).Return(
		nil,
		errors.New("test"),
	)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/matrix/channels/1/outputs/rooms",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"message":"Internal Server Error","status":"error"}`, string(body))
}

func TestAPI_ListOutputRoomsHandlerWithInvalidPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, _, server := testServer(ctrl)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/matrix/channels/0/outputs/rooms",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
