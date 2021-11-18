package handler

import (
	"fmt"
	"net/http"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
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
// @Tags Channels
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
// @Tags Channels
// @Security Admin-Authentication
// @Produce json
// @Param id path string true "Internal channel ID"
// @Success 200 {object} types.MessageSuccessResponse
// @Failure 401 {object} types.MessageErrorResponse
// @Router /channel/{id} [delete]
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

// PutUser godoc
// @Summary Change a User
// @Description Changes the settings or data for a matrix user.
// @Tags Users
// @Security Admin-Authentication
// @Accept json
// @Produce json
// @Param id path string true "Matrix account ID, user URL encoding where required"
// @Param blocked query boolean false "user state, if blocked no interaction with the bot is possible"
// @Param block_reason query string false "internally displayed reason for a block"
// @Success 200 {object} types.MessageSuccessResponse
// @Failure 401 {object} types.MessageErrorResponse
// @Router /user/{id} [put]
func (databaseHandler *DatabaseHandler) PutUser(ctx *gin.Context) {
	userID, err := getStringFromContext(ctx, "id")
	if err != nil {
		abort(ctx, http.StatusUnprocessableEntity, ResponseMessageNoID, err)
		return
	}

	data := &putUserData{}
	err = ctx.BindJSON(data)
	if err != nil {
		abort(ctx, http.StatusUnprocessableEntity, ResponseMessageNoID, err)
		return
	}
	log.Warn(fmt.Sprint(data))

	if data.Blocked != nil && *data.Blocked {
		channels, err := databaseHandler.database.GetChannelsByUserIdentifier(userID)
		if err != nil {
			abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
			return
		}

		for _, channel := range channels {
			err = databaseHandler.database.DeleteChannel(&channel)
			if err != nil {
				abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
				return
			}
		}

		err = databaseHandler.database.AddUserToBlocklist(userID, data.BlockReason)
		if err != nil {
			abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
			return
		}
	}

	response := types.MessageSuccessResponse{
		Status:  "success",
		Message: "Updated the user",
	}

	ctx.JSON(http.StatusOK, response)
}

type putUserData struct {
	Blocked     *bool  `json:"blocked"`
	BlockReason string `json:"block_reason"`
}
