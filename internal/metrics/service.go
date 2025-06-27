package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config for metrics service.
type Config struct {
	Address string
}

type service struct {
	server *http.Server
}

func New(cfg *Config) (Service, error) {
	s := http.Server{
		Addr:              cfg.Address,
		Handler:           promhttp.Handler(),
		ReadHeaderTimeout: time.Second * 30,
		WriteTimeout:      time.Second * 30,
	}

	return &service{
		server: &s,
	}, nil
}

func (service *service) Start() error {
	return service.server.ListenAndServe()
}

func (service *service) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	return service.server.Shutdown(ctx)
}
