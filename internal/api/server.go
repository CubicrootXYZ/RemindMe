package api

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/handler"
	"github.com/gin-gonic/gin"
)

// Server serves a http server
type Server struct {
	config   *configuration.Webserver
	calendar *handler.CalendarHandler
}

// NewServer returns a new webserver
func NewServer(config *configuration.Webserver, calendarHandler *handler.CalendarHandler) *Server {
	if len(config.APIkey) < 20 {
		panic(errors.ErrAPIkeyCriteriaNotMet)
	}
	return &Server{
		config:   config,
		calendar: calendarHandler,
	}
}

// Start starts the http server
func (server *Server) Start() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	calendarGroup := r.Group("/calendar")
	calendarGroup.Use(RequireAPIkey(server.config.APIkey))
	{
		calendarGroup.GET("", server.calendar.GetCalendars)
		calendarGroup.GET("/:id/ical", RequireIDInURI(), server.calendar.GetCalendarICal)
	}

	r.Run() // Port 8080
}
