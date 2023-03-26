package api

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/apictx"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/response"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/gin-gonic/gin"
)

type Room struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"` // RFC3339 formatted timestamp
	RoomID    string `json:"room_id"`    // Matrix room identifier
	Encrypted bool   `json:"encrypted"`  // True if an encrypted event is known
	Users     []User `json:"users"`      // List of matrix user identifiers known in this room
}

type User struct {
	ID      string `json:"id"`      // Matrix user identifier
	Blocked bool   `json:"blocked"` // If true the user can not interact with the bot
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

// @Summary List all Input Rooms
// @Description List all matrix rooms acting as an input for the given channel.
// @Tags Matrix
// @Security APIKeyAuthentication
// @Produce json
// @Param id path string true "Channel ID"
// @Success 200 {object} response.DataResponse{data=[]Room}
// @Failure 401 {object} response.MessageErrorResponse
// @Failure 404 {object} response.MessageErrorResponse
// @Failure 500 ""
// @Router /matrix/channels/{id}/inputs/rooms [get]
func (api *api) listInputRoomsHandler(ctx *gin.Context) {
	channelID, ok := apictx.GetUintFromContext(ctx, "id")
	if !ok || channelID < 1 {
		response.AbortWithNotFoundError(ctx)
		return
	}

	rooms, err := api.config.MatrixDB.ListInputRoomsByChannel(channelID)
	if err != nil {
		if err != nil {
			api.logger.Err(err)
			response.AbortWithInternalServerError(ctx)
			return
		}
	}

	response.WithData(ctx, roomsToResponse(rooms))
}

// @Summary List all Output Rooms
// @Description List all matrix rooms acting as an output for the given channel.
// @Tags Matrix
// @Security APIKeyAuthentication
// @Produce json
// @Param id path string true "Channel ID"
// @Success 200 {object} response.DataResponse{data=[]Room}
// @Failure 401 {object} response.MessageErrorResponse
// @Failure 404 {object} response.MessageErrorResponse
// @Failure 500 ""
// @Router /matrix/channels/{id}/outputs/rooms [get]
func (api *api) listOutputRoomsHandler(ctx *gin.Context) {
	channelID, ok := apictx.GetUintFromContext(ctx, "id")
	if !ok || channelID < 1 {
		response.AbortWithNotFoundError(ctx)
		return
	}

	rooms, err := api.config.MatrixDB.ListOutputRoomsByChannel(channelID)
	if err != nil {
		if err != nil {
			api.logger.Err(err)
			response.AbortWithInternalServerError(ctx)
			return
		}
	}

	response.WithData(ctx, roomsToResponse(rooms))
}
