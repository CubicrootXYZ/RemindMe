package database_test

import (
	"fmt"
	"math/rand"
	"testing"

	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testUser() *matrixdb.MatrixUser {
	room, err := service.NewRoom(testRoom())
	if err != nil {
		panic(err)
	}

	userID := "@remindme:example.org"
	for err == nil {
		userID = fmt.Sprintf("@%d:example.org", rand.Int()) //nolint:gosec
		_, err = service.GetUserByID(userID)
	}

	return &matrixdb.MatrixUser{
		ID:    userID,
		Rooms: []matrixdb.MatrixRoom{*room},
	}
}

func TestService_GetUserByID(t *testing.T) {
	userBefore, err := service.NewUser(testUser())
	require.NoError(t, err)

	userAfter, err := service.GetUserByID(userBefore.ID)
	require.NoError(t, err)

	assertUsersEqual(t, userBefore, userAfter)
}

func TestService_GetUserByIDWithUserNotFound(t *testing.T) {
	_, err := service.GetUserByID("abc")
	assert.ErrorIs(t, err, matrixdb.ErrNotFound)
}

func assertUsersEqual(t *testing.T, a *matrixdb.MatrixUser, b *matrixdb.MatrixUser) {
	t.Helper()

	assert.Equal(t, a.ID, b.ID)

	require.Equal(t, len(a.Rooms), len(b.Rooms))
	for i := range a.Rooms {
		assertRoomsEqual(t, &a.Rooms[i], &b.Rooms[i])
	}
}
