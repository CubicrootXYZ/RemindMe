package api

import (
	"log/slog"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/apictx"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
)

type api struct {
	config *Config
	logger *slog.Logger
}

// Config holds the configuration for the API.
type Config struct {
	Database database.Service
	MatrixDB matrixdb.Service

	DefaultAuthProvider gin.HandlerFunc
}

// New assembles a new API.
func New(config *Config, logger *slog.Logger) API {
	return &api{
		config: config,
		logger: logger,
	}
}

func (api *api) RegisterRoutes(r *gin.Engine) error {
	router := r.Group("/matrix")

	channels := router.Group("/channels")
	channels.Use(api.config.DefaultAuthProvider)
	channels.GET("/:id/inputs/rooms", apictx.RequireIDInURI(), api.listInputRoomsHandler)
	channels.GET("/:id/outputs/rooms", apictx.RequireIDInURI(), api.listOutputRoomsHandler)

	return nil
}
