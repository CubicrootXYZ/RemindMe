package api

import (
	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/apictx"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/gin-gonic/gin"
)

// Config holds information for the API service.
type Config struct {
	IcalDB icaldb.Service
}

type api struct {
	icalDB icaldb.Service
	logger gologger.Logger
}

// New assembles a new iCal API.
func New(config *Config, logger gologger.Logger) API {
	return &api{
		icalDB: config.IcalDB,
		logger: logger,
	}
}

func (api *api) RegisterRoutes(r *gin.Engine) error {
	router := r.Group("/ical")
	router.GET("/:id", apictx.RequireIDInURI(), api.icalExportHandler)

	return nil
}
