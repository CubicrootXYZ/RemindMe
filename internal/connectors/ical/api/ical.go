package api

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/apictx"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/response"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
)

// GetCalendarICal godoc
// @Summary Get calendar (iCal)
// @Description Get calendar as iCal
// @Tags Calendars
// @Produce plain
// @Param id path int true "iCal output ID"
// @Param token query string true "authentication token"
// @Success 200 ""
// @Failure 400 {object} response.MessageErrorResponse
// @Failure 401 {object} response.MessageErrorResponse
// @Failure 500 ""
// @Router /ical/{id} [get]
func (api *api) icalExportHandler(ctx *gin.Context) {
	id, ok := apictx.GetUintFromContext(ctx, "id")
	if !ok {
		response.AbortWithNotFoundError(ctx)
		return
	}
	token := ctx.Query("token")
	if token == "" {
		// Use not found to not leak any information.
		response.AbortWithNotFoundError(ctx)
		return
	}

	output, err := api.icalDB.GetIcalOutputByID(id)
	if err != nil {
		if errors.Is(err, icaldb.ErrNotFound) {
			response.AbortWithNotFoundError(ctx)
			return
		}
		api.logger.Error("failed to get iCal output", "error", err, "ical.output.id", id)
		response.AbortWithInternalServerError(ctx)
		return
	}

	if i := subtle.ConstantTimeCompare([]byte(output.Token), []byte(token)); i != 1 {
		response.AbortWithNotFoundError(ctx)
		return
	}

	o, err := api.database.GetOutputByType(output.ID, ical.OutputType)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			response.AbortWithNotFoundError(ctx)
			return
		}
		api.logger.Error("failed to get output from database", "error", err)
		response.AbortWithInternalServerError(ctx)
		return
	}

	events, err := api.database.GetEventsByChannel(o.ChannelID)
	if err != nil {
		api.logger.Error("failed to get events from database", "error", err)
		response.AbortWithInternalServerError(ctx)
		return
	}

	calendar := format.NewCalendar(strconv.Itoa(int(o.ChannelID)), events)

	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%d.ics\"", o.ChannelID))
	ctx.String(http.StatusOK, calendar)
}
