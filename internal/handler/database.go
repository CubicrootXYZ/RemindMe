package handler

import (
	"net/http"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
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
		channelsPublic = append(channelsPublic, channelToChannelResponse(&channel))
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
// @Param blocked body boolean false "user state, if blocked no interaction with the bot is possible"
// @Param block_reason body string false "internally displayed reason for a block"
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
	} else if data.Blocked != nil && !*data.Blocked {
		err = databaseHandler.database.RemoveUserFromBlocklist(userID)
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

// GetUsers godoc
// @Summary Get all Users
// @Description Lists all users and their channels
// @Tags Users
// @Security Admin-Authentication
// @Produce json
// @Param include[] query string false "Comma separated list of additional users to include. One of: blocked" collectionFormat(multi)
// @Success 200 {object} types.DataResponse{data=[]userResponse}
// @Failure 401 {object} types.MessageErrorResponse
// @Router /user [get]
func (databaseHandler *DatabaseHandler) GetUsers(ctx *gin.Context) {
	data := &getUsersData{}
	err := ctx.Bind(data)
	if err != nil {
		abort(ctx, http.StatusUnprocessableEntity, ResponseMessageNoID, err)
		return
	}

	channels, err := databaseHandler.database.GetChannelList()
	if err != nil {
		abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
		return
	}

	users := channelsToUserList(channels)

	for _, group := range uniqueString(data.Include) {
		switch group {
		case "blocked":
			blocklists, err := databaseHandler.database.GetBlockedUserList()
			if err != nil {
				abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
				return
			}

			users = append(users, blocklistsToUserList(blocklists)...)
		}
	}

	response := types.DataResponse{
		Status: "success",
		Data:   users,
	}

	ctx.JSON(http.StatusOK, response)
}

// Helper

func channelsToUserList(channels []database.Channel) []*userResponse {
	responseData := make([]*userResponse, 0)

CHANNELS:
	for _, channel := range channels {
		for _, user := range responseData {
			if user.UserIdentifier == channel.UserIdentifier {
				user.Channels = append(user.Channels, channelToChannelResponse(&channel))
				continue CHANNELS
			}
		}

		responseData = append(responseData, &userResponse{
			UserIdentifier: channel.UserIdentifier,
			Blocked:        false,
			Channels:       []channelResponse{channelToChannelResponse(&channel)},
		})
	}

	return responseData
}

func blocklistsToUserList(blocklists []database.Blocklist) []*userResponse {
	responseData := make([]*userResponse, 0)
	for _, blocklist := range blocklists {
		responseData = append(responseData, &userResponse{
			UserIdentifier: blocklist.UserIdentifier,
			Blocked:        true,
			Comment:        blocklist.Reason,
			Channels:       []channelResponse{},
		})
	}

	return responseData
}

func channelToChannelResponse(channel *database.Channel) channelResponse {
	return channelResponse{
		ID:                channel.ID,
		Created:           channel.Created,
		ChannelIdentifier: channel.ChannelIdentifier,
		UserIdentifier:    channel.UserIdentifier,
		TimeZone:          channel.TimeZone,
		DailyReminder:     channel.DailyReminder == nil,
		Role:              channel.Role,
	}
}

func uniqueString(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
