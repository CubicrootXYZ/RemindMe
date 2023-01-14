package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/handler"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// Server serves a http server
type Server struct {
	config   *configuration.Webserver
	calendar *handler.CalendarHandler
	database *handler.DatabaseHandler
	server   *http.Server
}

type Handler struct {
	Calendar *handler.CalendarHandler
	Database *handler.DatabaseHandler
}

// NewServer returns a new webserver
func NewServer(config *configuration.Webserver, handler *Handler) *Server {
	if len(config.APIkey) < 20 {
		panic(errors.ErrAPIkeyCriteriaNotMet)
	}
	return &Server{
		config:   config,
		calendar: handler.Calendar,
		database: handler.Database,
	}
}

// Start starts the http server
func (server *Server) Start(debug bool) error {
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
		channelGroup.POST("/:id/thirdpartyresources", RequireIDInURI(), server.database.PostChannelThirdPartyResource)
		channelGroup.GET("/:id/thirdpartyresources", RequireIDInURI(), server.database.GetChannelThirdPartyResource)
		channelGroup.DELETE("/:id/thirdpartyresources/:id2", RequireIDInURI(), RequireID2InURI(), server.database.DeleteChannelThirdPartyResource)
	}

	userGroup := r.Group("/user")
	userGroup.Use(RequireAPIkey(server.config.APIkey))
	{
		userGroup.GET("", server.database.GetUsers)
		userGroup.PUT("/:id", RequireStringIDInURI(), server.database.PutUser)
	}

	r.GET("calendar/:id/ical", RequireCalendarSecret(), RequireIDInURI(), server.calendar.GetCalendarICal)

	server.server = &http.Server{
		Addr:         server.config.Address,
		Handler:      r,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}
	if err := server.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error(fmt.Sprintf("Error when starting server: %s", err.Error()))
		return err
	}
	log.Info("server stopped")
	return nil
}

func (server *Server) Stop(ctx context.Context) error {
	log.Debug("stopping server ...")
	return server.server.Shutdown(ctx)
}
