package coreapi

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/response"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
)

type Channel struct {
	ID            uint
	CreatedAt     string // RFC 3339 formated time
	Description   string
	DailyReminder *string // HH:MM of daily reminder or null if disabled
}

func channelToResponse(channelIn *database.Channel) Channel {
	channelOut := Channel{
		ID:          channelIn.ID,
		CreatedAt:   channelIn.CreatedAt.Format(time.RFC3339),
		Description: channelIn.Description,
	}

	if channelIn.DailyReminder != nil {
		dailyReminder := fmt.Sprintf("%02d:%02d", int(*channelIn.DailyReminder/60), *channelIn.DailyReminder%60)
		channelOut.DailyReminder = &dailyReminder
	}

	return channelOut
}

func channelsToResponse(channelsIn []database.Channel) []Channel {
	channelsOut := make([]Channel, len(channelsIn))

	for i := range channelsIn {
		channelsOut[i] = channelToResponse(&channelsIn[i])
	}

	return channelsOut
}

// listChannelsHandler godoc
// @Summary List all channels
// @Description List all channels
// @Tags Channels
// @Security APIKeyAuthentication
// @Produce json
// @Success 200 {object} response.DataResponse{data=[]Channel}
// @Failure 401 {object} response.MessageErrorResponse
// @Failure 500 ""
// @Router /core/channels [get]
func (api *coreAPI) listChannelsHandler(ctx *gin.Context) {
	channels, err := api.config.Database.GetChannels()
	if err != nil {
		api.logger.Error("failed to list channels", "error", err)
		response.AbortWithInternalServerError(ctx)
		return
	}

	response.WithData(ctx, channelsToResponse(channels))
}
