package api

import (
	"time"

	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/gin-gonic/gin"
)

type Room struct {
	ID        uint
	CreatedAt string // RFC3339 formatted timestamp
	RoomID    string // Matrix room identifier
	Encrypted bool   // True if an encrypted event is known
	Users     []User // List of matrix user identifiers known in this room
}

type User struct {
	ID      string // Matrix user identifier
	Blocked bool   // If true the user can not interact with the bot
}

func roomToResponse(room *matrixdb.MatrixRoom) Room {
	users := make([]User, len(room.Users))

	for i := range room.Users {
		users[i] = User{
			ID:      room.Users[i].ID,
			Blocked: room.Users[i].Blocked,
		}
	}
	return Room{
		ID:        room.ID,
		CreatedAt: room.CreatedAt.Format(time.RFC3339),
		RoomID:    room.RoomID,
		Encrypted: room.LastCryptoEvent != "",
		Users:     users,
	}
}

func roomsToResponse(rooms []matrixdb.MatrixRoom) []Room {
	roomsOut := make([]Room, len(rooms))

	for i := range rooms {
		roomsOut[i] = roomToResponse(&rooms[i])
	}

	return roomsOut
}

// listRoomsHandler godoc
// @Summary List all Rooms
// @Description List all matrix rooms.
// @Tags Matrix
// @Security APIKeyAuthentication
// @Produce json
// @Param id path string true "Channel ID"
// @Success 200 {object} response.DataResponse{data=[]Room}
// @Failure 401 {object} response.MessageErrorResponse
// @Failure 500 ""
// @Router /matrix/channels/{id}/rooms [get]
func (api *api) listRoomsHandler(ctx *gin.Context) {
	// TODO test
}
