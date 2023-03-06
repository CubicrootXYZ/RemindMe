package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type server struct {
	server *http.Server
	config *Config
	logger gologger.Logger
}

type EndpointProvider interface {
	RegisterRoutes(*gin.Engine) error
}

// Config for the server.
type Config struct {
	EndpointProviders map[string]EndpointProvider
	Address           string
}

// NewServer assembles a new API webserver.
func NewServer(config *Config, logger gologger.Logger) Server {
	return &server{
		config: config,
		logger: logger,
	}
}

// Start the webserver and serve the given endpoints.
// Blocks until stopped.
func (server *server) Start() error {
	server.logger.Infof("starting server at '%s'", server.config.Address)
	err := server.assembleRoutes()
	if err != nil {
		return err
	}

	if server.server == nil {
		return errors.New("server setup failed, can not start with empty server")
	}

	if err := server.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop the server.
// Might take a few moments.
func (server *server) Stop() error {
	server.logger.Infof("stopping server ...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return server.server.Shutdown(ctx)
}
