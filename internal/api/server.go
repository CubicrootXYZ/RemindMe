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
	database *handler.DatabaseHandler
}

// NewServer returns a new webserver
func NewServer(config *configuration.Webserver, calendarHandler *handler.CalendarHandler, databaseHandler *handler.DatabaseHandler) *Server {
	if len(config.APIkey) < 20 {
		panic(errors.ErrAPIkeyCriteriaNotMet)
	}
	return &Server{
		config:   config,
		calendar: calendarHandler,
		database: databaseHandler,
	}
}

// Start starts the http server
func (server *Server) Start(debug bool) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Logger())
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	calendarGroup := r.Group("/calendar")
	calendarGroup.Use(RequireAPIkey(server.config.APIkey))
	{
		calendarGroup.GET("", server.calendar.GetCalendars)
		calendarGroup.PATCH("/:id", RequireIDInURI(), server.calendar.PatchCalender)
	}

	channelGroup := r.Group("/channel")
	channelGroup.Use(RequireAPIkey(server.config.APIkey))
	{
		channelGroup.GET("", server.database.GetChannels)
		channelGroup.DELETE("/:id", RequireIDInURI(), server.database.DeleteChannel)
	}

	userGroup := r.Group("/user")
	userGroup.Use(RequireAPIkey(server.config.APIkey))
	{
		userGroup.GET("", server.database.GetUsers)
		userGroup.PUT("/:id", RequireStringIDInURI(), server.database.PutUser)
	}

	r.GET("calendar/:id/ical", RequireCalendarSecret(), RequireIDInURI(), server.calendar.GetCalendarICal)

	r.Run() // Port 8080
}
