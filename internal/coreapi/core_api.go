package coreapi

import (
	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
)

type coreAPI struct {
	config *Config
	logger gologger.Logger
}

// Config holds the configuration for the core API.
type Config struct {
	Database            database.Service
	DefaultAuthProvider gin.HandlerFunc
}

// New assembles a new core API.
// Core API uses the /core path prefix.
func New(config *Config, logger gologger.Logger) CoreAPI {
	return &coreAPI{
		config: config,
		logger: logger,
	}
}

// RegisterRoutes registers the routes for the core API.
func (api *coreAPI) RegisterRoutes(r *gin.Engine) error {
	router := r.Group("/core")

	router.GET("/channels", api.config.DefaultAuthProvider, api.listChannelsHandler)

	return nil
}
