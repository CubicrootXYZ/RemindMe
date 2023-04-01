package api

import (
	"crypto/subtle"
	"errors"
	"fmt"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/apictx"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/response"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
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
// @Failure 400 {object} types.MessageErrorResponse
// @Failure 401 {object} types.MessageErrorResponse
// @Failure 500 ""
// @Router /ical/{id} [get]
func (api *api) icalExportHandler(ctx *gin.Context) {
	// TODO test
	id, ok := apictx.GetUintFromContext(ctx, "id")
	if !ok {
		response.AbortWithNotFoundError(ctx)
		return
	}
	token, ok := ctx.Params.Get("token")
	if !ok {
		// Use not found to not leak any information.
		response.AbortWithNotFoundError(ctx)
		return
	}

	output, err := api.icalDB.GetIcalOutputByID(id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			response.AbortWithNotFoundError(ctx)
			return
		}
		response.AbortWithInternalServerError(ctx)
		return
	}

	if i := subtle.ConstantTimeCompare([]byte(output.Token), []byte(token)); i != 1 {
		response.AbortWithNotFoundError(ctx)
		return
	}

	// TODO get channel and events
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%d.ics\"", "TODO"))
}
