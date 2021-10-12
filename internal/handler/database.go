package handler

import (
	"net/http"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"github.com/gin-gonic/gin"
)

// DatabaseHandler groups the handlers for the calendars
type DatabaseHandler struct {
	database types.Database
}

// NewDatabaseHandler makes a new handler
func NewDatabaseHandler(database types.Database) *DatabaseHandler {
	return &DatabaseHandler{
		database: database,
	}
}

// GetChannels godoc
// @Summary List all channels
// @Description List all channels
// @Security Admin-Authentication
// @Produce json
// @Success 200 {object} types.DataResponse{data=[]channelResponse}
// @Failure 401 {object} types.MessageErrorResponse
// @Router /channel [get]
func (databaseHandler *DatabaseHandler) GetChannels(ctx *gin.Context) {
	channelsPublic := make([]channelResponse, 0)

	channels, err := databaseHandler.database.GetChannelList()
	if err != nil {
		abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
		return
	}

	for _, channel := range channels {
		channelsPublic = append(channelsPublic, channelResponse{
			ID:                channel.ID,
			Created:           channel.Created,
			ChannelIdentifier: channel.ChannelIdentifier,
			UserIdentifier:    channel.UserIdentifier,
			TimeZone:          channel.TimeZone,
			DailyReminder:     channel.DailyReminder == nil,
			Role:              channel.Role,
		})
	}

	response := types.DataResponse{
		Status: "success",
		Data:   channelsPublic,
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteChannel godoc
// @Summary Delete a channel
// @Description Delete a channel and remove access for this user. If the bot is open for invites the user can simply start a new chat.
// @Security Admin-Authentication
// @Produce json
// @Success 200 {object} types.MessageSuccessResponse
// @Failure 401 {object} types.MessageErrorResponse
// @Router /channel/{id} [get]
func (databaseHandler *DatabaseHandler) DeleteChannel(ctx *gin.Context) {
	channelID, err := getUintFromContext(ctx, "id")
	if err != nil {
		abort(ctx, http.StatusUnprocessableEntity, ResponseMessageNoID, err)
		return
	}

	channel, err := databaseHandler.database.GetChannel(channelID)
	if err != nil {
		abort(ctx, http.StatusNotFound, ResponseMessageNotFound, err)
		return
	}

	err = databaseHandler.database.DeleteChannel(channel)
	if err != nil {
		abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
		return
	}

	response := types.MessageSuccessResponse{
		Status:  "success",
		Message: "Deleted the channel",
	}

	ctx.JSON(http.StatusOK, response)
}
