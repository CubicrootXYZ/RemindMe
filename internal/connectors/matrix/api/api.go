package api

import (
	"github.com/CubicrootXYZ/gologger"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
)

type api struct {
	config *Config
	logger gologger.Logger
}

// Config holds the configuration for the API.
type Config struct {
	Database database.Service
	MatrixDB matrixdb.Service
}

// New assembles a new API.
func New(config *Config, logger gologger.Logger) API {
	return &api{
		config: config,
		logger: logger,
	}
}

func (api *api) RegisterRoutes(r *gin.Engine) error {
	_ = r.Group("matrix")

	// TODO

	return nil
}
