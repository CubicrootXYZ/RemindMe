package handler

import (
	"fmt"
	"net/http"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/calendar"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"github.com/gin-gonic/gin"
)

// CalendarHandler groups the handlers for the calendars
type CalendarHandler struct {
	database types.Database
}

// NewCalendarHandler makes a new handler
func NewCalendarHandler(database types.Database) *CalendarHandler {
	return &CalendarHandler{
		database: database,
	}
}

// GetCalendars godoc
// @Summary List all calendars
// @Description List all available calendars
// @Accept query
// @Produce json
// @Success 200 {object} types.DataResponse{Data=[]calendar}
// @Failure 401 {object} types.Unauthenticated
// @Router /calendar [get]
func (calendarHandler *CalendarHandler) GetCalendars(ctx *gin.Context) {
	calendars := make([]calendarResponse, 0)

	channels, err := calendarHandler.database.GetChannelList()
	if err != nil {
		abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
		return
	}

	for _, channel := range channels {
		calendars = append(calendars, calendarResponse{
			ID:      channel.ID,
			User:    channel.UserIdentifier,
			Token:   channel.CalendarSecret,
			Channel: channel.ChannelIdentifier,
		})
	}

	response := types.DataResponse{
		Status: "success",
		Data:   calendars,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetCalendarICal godoc
// @Summary Get calendar (iCal)
// @Description Get calendar as iCal
// @Accept query
// @Produce plain
// @Param id path int true "Calendar ID"
// @Param token query string true "authentication token"
// @Success 200 {string}
// @Failure 401 {object} types.Unauthenticated
// @Router /calendar/{id}/ical [get]
func (calendarHandler *CalendarHandler) GetCalendarICal(ctx *gin.Context) {
	token, err := getStringFromContext(ctx, "token")
	if err != nil {
		abort(ctx, http.StatusUnauthorized, ResponseMessageUnauthorized, err)
		return
	}

	channelID, err := getUintFromContext(ctx, "id")
	if err != nil {
		abort(ctx, http.StatusUnauthorized, ResponseMessageUnauthorized, err)
		return
	}

	channel, err := calendarHandler.database.GetChannel(channelID)
	if err != nil {
		abort(ctx, http.StatusUnauthorized, ResponseMessageUnauthorized, err)
		return
	}

	if channel.CalendarSecret != token || len(channel.CalendarSecret) < 20 {
		abort(ctx, http.StatusUnauthorized, ResponseMessageUnauthorized, err)
		return
	}

	reminders, err := calendarHandler.database.GetPendingReminders(channel)
	if err != nil {
		abort(ctx, http.StatusInternalServerError, ResponseMessageInternalServerError, err)
		return
	}

	calendar := calendar.NewCalendar(&reminders)
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%d.ics\"", channelID))
	ctx.String(http.StatusOK, calendar.ICal())
}
