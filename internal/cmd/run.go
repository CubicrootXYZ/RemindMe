package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"golang.org/x/sync/errgroup"
)

// Run setups the application and runs it
func Run(config *Config) error {
	logger := gologger.New(config.loggerConfig(), 0).WithField("component", "cmd")
	defer logger.Flush()

	logger.Infof("starting up RemindMe ...")
	processes, err := setup(config, logger)
	if err != nil {
		logger.Err(err)
		return err
	}

	eg, ctx := errgroup.WithContext(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		select {
		case s := <-sigChan:
			logger.Infof("received signal, shutting down: %v", s)
		case <-ctx.Done():
			logger.Infof("at least one process exited, shutting down")
		}

		for _, p := range processes {
			err := p.Stop()
			if err != nil {
				logger.Err(err)
			}
		}
	}()

	for _, p := range processes {
		p := p
		eg.Go(func() error {
			return p.Start()
		})
	}

	err = eg.Wait()
	if err != nil {
		logger.Err(err)
	}

	logger.Infof("shut down complete, bye")
	return err
}

type process interface {
	Start() error
	Stop() error
}

func setup(config *Config, logger gologger.Logger) ([]process, error) {
	db, err := database.NewService(config.databaseConfig(), logger.WithField("component", "database"))
	if err != nil {
		logger.Err(err)
		return nil, err
	}

	daemon := daemon.New(config.daemonConfig(), db, logger.WithField("component", "daemon"))

	return []process{daemon}, nil
}
