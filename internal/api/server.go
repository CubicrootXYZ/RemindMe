package api

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"context"

	"github.com/gin-gonic/gin"
)

type server struct {
	server *http.Server
	config *Config
	logger *slog.Logger
}

type RouteProvider interface {
	RegisterRoutes(*gin.Engine) error
}

// Config for the server.
type Config struct {
	RouteProviders map[string]RouteProvider
	Address        string
}

// NewServer assembles a new API webserver.
func NewServer(config *Config, logger *slog.Logger) Server {
	return &server{
		config: config,
		logger: logger,
	}
}

// Start the webserver and serve the given endpoints.
// Blocks until stopped.
func (server *server) Start() error {
	server.logger.Info("starting server", "address", server.config.Address)

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
	timeout := time.Second * 5
	server.logger.Info("stopping server", "timeout", timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return server.server.Shutdown(ctx)
}
