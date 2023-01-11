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
// @Security AdminAuthentication
// @Produce json
// @Success 200 {object} types.DataResponse{data=[]channelResponse}
// @Failure 401 {object} types.MessageErrorResponse
// @Failure 500 ""
// @Router /channel [get]
func (databaseHandler *DatabaseHandler) GetChannels(ctx *gin.Context) {
	channelsPublic := make([]channelResponse, 0)

	channels, err := databaseHandler.database.GetChannelList()
	if err != nil {
		abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
		return
	}

	for i := range channels {
		channelsPublic = append(channelsPublic, channelToChannelResponse(&channels[i]))
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
// @Security AdminAuthentication
// @Produce json
// @Param id path string true "Internal channel ID"
// @Success 200 {object} types.MessageSuccessResponse
// @Failure 401 {object} types.MessageErrorResponse
// @Failure 404 {object} types.MessageErrorResponse
// @Failure 422 {object} types.MessageErrorResponse "Input validation failed"
// @Failure 500 ""
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

// GetChannelThirdPartyResource godoc
// @Summary Get third party resources
// @Description Lists all third party resources in this channel.
// @Tags Channels
// @Security AdminAuthentication
// @Produce json
// @Param id path string true "Internal channel ID"
// @Success 200 {object} types.DataResponse{data=[]thirdPartyResourceResponse}
// @Failure 401 {object} types.MessageErrorResponse
// @Failure 404 {object} types.MessageErrorResponse
// @Failure 422 {object} types.MessageErrorResponse "Input validation failed"
// @Failure 500 ""
// @Router /channel/{id}/thirdpartyresources [get]
func (databaseHandler *DatabaseHandler) GetChannelThirdPartyResource(ctx *gin.Context) {
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

	resources, err := databaseHandler.database.GetThirdPartyResourcesByChannel(channel.ID)
	if err != nil {
		abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
		return
	}

	response := types.DataResponse{
		Status: "success",
		Data:   thirdPartyResourcesToResponse(resources),
	}

	ctx.JSON(http.StatusOK, response)
}

// PostChannelThirdPartyResource godoc
// @Summary Add a third party resource to a channel
// @Description Add a third party resource to a channel.
// @Tags Channels
// @Security AdminAuthentication
// @Produce json
// @Param id path string true "Internal channel ID"
// @Param payload body postChannelThirdPartyResourceData true "payload"
// @Success 200 {object} types.MessageSuccessResponse
// @Failure 401 {object} types.MessageErrorResponse
// @Failure 404 {object} types.MessageErrorResponse
// @Failure 422 {object} types.MessageErrorResponse "Input validation failed"
// @Failure 500 ""
// @Router /channel/{id}/thirdpartyresources [post]
func (databaseHandler *DatabaseHandler) PostChannelThirdPartyResource(ctx *gin.Context) {
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

	data := &postChannelThirdPartyResourceData{}
	err = ctx.ShouldBindJSON(data)
	if err != nil {
		abort(ctx, http.StatusUnprocessableEntity, ResponseMessageInvalidData, err)
		return
	}

	resourceType, err := database.ThirdPartyResourceTypeFromString(data.Type)
	if err != nil {
		abort(ctx, http.StatusUnprocessableEntity, ResponseMessageUnknownType, err)
		return
	}

	if data.ResourceURL == "" {
		abort(ctx, http.StatusUnprocessableEntity, ResponseMessageMissingURL, nil)
	}

	_, err = databaseHandler.database.AddThirdPartyResource(&database.ThirdPartyResource{
		Type:        resourceType,
		ChannelID:   channel.ID,
		ResourceURL: data.ResourceURL,
	})
	if err != nil {
		abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
		return
	}

	response := types.MessageSuccessResponse{
		Status:  "success",
		Message: "Added the resource",
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteChannelThirdPartyResource godoc
// @Summary Delete a third party resource
// @Description Delete a third party resource.
// @Tags Channels
// @Security AdminAuthentication
// @Produce json
// @Param id path string true "Internal channel ID"
// @Param id2 path string true "Internal third party resource ID"
// @Param payload body postChannelThirdPartyResourceData true "payload"
// @Success 200 {object} types.MessageSuccessResponse
// @Failure 401 {object} types.MessageErrorResponse
// @Failure 404 {object} types.MessageErrorResponse
// @Failure 422 {object} types.MessageErrorResponse "Input validation failed"
// @Failure 500 ""
// @Router /channel/{id}/thirdpartyresources/{id2} [delete]
func (databaseHandler *DatabaseHandler) DeleteChannelThirdPartyResource(ctx *gin.Context) {
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

	resourceID, err := getUintFromContext(ctx, "id2")
	if err != nil {
		abort(ctx, http.StatusUnprocessableEntity, ResponseMessageNoID, err)
		return
	}

	resources, err := databaseHandler.database.GetThirdPartyResourcesByChannel(channel.ID)
	if err != nil {
		abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
		return
	}

	resourceIsInChannel := false
	for _, resource := range resources {
		if resource.ID == resourceID {
			resourceIsInChannel = true
			break
		}
	}

	if !resourceIsInChannel {
		abort(ctx, http.StatusNotFound, ResponseMessageNotFound, nil)
		return
	}

	err = databaseHandler.database.DeleteThirdPartyResource(resourceID)
	if err != nil {
		abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
		return
	}

	response := types.MessageSuccessResponse{
		Status:  "success",
		Message: "Deleted the resource",
	}

	ctx.JSON(http.StatusOK, response)
}

// PutUser godoc
// @Summary Change a User
// @Description Changes the settings or data for a matrix user.
// @Tags Users
// @Security AdminAuthentication
// @Accept json
// @Produce json
// @Param id path string true "Matrix account ID, use URL encoding"
// @Param payload body putUserData true "payload"
// @Success 200 {object} types.MessageSuccessResponse
// @Failure 401 {object} types.MessageErrorResponse
// @Failure 422 {object} types.MessageErrorResponse "Input validation failed"
// @Failure 500 ""
// @Router /user/{id} [put]
func (databaseHandler *DatabaseHandler) PutUser(ctx *gin.Context) {
	userID, err := getStringFromContext(ctx, "id")
	if err != nil {
		abort(ctx, http.StatusUnprocessableEntity, ResponseMessageNoID, err)
		return
	}

	data := &putUserData{}
	err = ctx.ShouldBindJSON(data)
	if err != nil {
		abort(ctx, http.StatusUnprocessableEntity, ResponseMessageNoID, err)
		return
	}

	if data.Blocked != nil && *data.Blocked {
		err := databaseHandler.blockUserID(userID, data.BlockReason)
		if err != nil {
			abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
			return
		}
	} else if data.Blocked != nil && !*data.Blocked {
		err = databaseHandler.unblockUserID(userID)
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

func (databaseHandler *DatabaseHandler) blockUserID(userID string, blockReason string) error {
	channels, err := databaseHandler.database.GetChannelsByUserIdentifier(userID)
	if err != nil {
		return err
	}

	for i := range channels {
		err = databaseHandler.database.DeleteChannel(&channels[i])
		if err != nil {
			return err
		}
	}

	err = databaseHandler.database.AddUserToBlocklist(userID, blockReason)
	if err != nil {
		return err
	}

	return nil
}

func (databaseHandler *DatabaseHandler) unblockUserID(userID string) error {
	return databaseHandler.database.RemoveUserFromBlocklist(userID)
}

// GetUsers godoc
// @Summary Get all Users
// @Description Lists all users and their channels
// @Tags Users
// @Security AdminAuthentication
// @Produce json
// @Param include[] query []string false "Comma separated list of additional users to include. One of: blocked" collectionFormat(multi)
// @Success 200 {object} types.DataResponse{data=[]userResponse}
// @Failure 401 {object} types.MessageErrorResponse
// @Failure 422 {object} types.MessageErrorResponse "Input validation failed"
// @Failure 500 ""
// @Router /user [get]
func (databaseHandler *DatabaseHandler) GetUsers(ctx *gin.Context) {
	data := &getUsersData{}
	err := ctx.ShouldBind(data)
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
	for i := range channels {
		for j := range responseData {
			if responseData[j].UserIdentifier == channels[i].UserIdentifier {
				responseData[j].Channels = append(responseData[j].Channels, channelToChannelResponse(&channels[i]))
				continue CHANNELS
			}
		}

		responseData = append(responseData, &userResponse{
			UserIdentifier: channels[i].UserIdentifier,
			Blocked:        false,
			Channels:       []channelResponse{channelToChannelResponse(&channels[i])},
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

func thirdPartyResourceToResponse(resource *database.ThirdPartyResource) thirdPartyResourceResponse {
	return thirdPartyResourceResponse{
		ID:          resource.ID,
		Type:        resource.Type.String(),
		ResourceURL: resource.ResourceURL,
	}
}

func thirdPartyResourcesToResponse(resources []database.ThirdPartyResource) []thirdPartyResourceResponse {
	response := make([]thirdPartyResourceResponse, len(resources))

	for i := range resources {
		response[i] = thirdPartyResourceToResponse(&resources[i])
	}

	return response
}
