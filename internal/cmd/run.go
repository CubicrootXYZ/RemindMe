package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/middleware"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	icalapi "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/api"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/reply"
	matrixapi "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/api"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/coreapi"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/rs/zerolog/log"
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
	processes := []process{}

	// Database
	dbConfig := config.databaseConfig()
	dbConfig.InputServices = make(map[string]database.InputService)
	dbConfig.OutputServices = make(map[string]database.OutputService)
	db, err := database.NewService(dbConfig, logger.WithField("component", "database"))
	if err != nil {
		logger.Err(err)
		return nil, err
	}

	// iCal connector
	icalDB, err := icaldb.New(db.GormDB())
	if err != nil {
		log.Err(err)
		return nil, err
	}
	icalConnector := ical.New(&ical.Config{
		ICalDB:   icalDB,
		Database: db,
	}, logger.WithField("component", "ical connector"))

	dbConfig.OutputServices[ical.OutputType] = icalConnector
	dbConfig.OutputServices[ical.InputType] = icalConnector

	// Matrix connector
	matrixDB, err := matrixdb.New(db.GormDB())
	if err != nil {
		logger.Err(err)
		return nil, err
	}

	matrixConnector, err := matrix.New(assembleMatrixConfig(config, icalConnector), db, matrixDB, logger.WithField("component", "matrix connector"))
	if err != nil {
		log.Err(err)
		return nil, err
	}
	processes = append(processes, matrixConnector)

	dbConfig.InputServices[matrix.InputType] = matrixConnector
	dbConfig.OutputServices[matrix.OutputType] = matrixConnector

	// Daemon
	daemonConf := config.daemonConfig()
	daemonConf.OutputServices = make(map[string]daemon.OutputService)
	daemonConf.OutputServices[matrix.OutputType] = matrixConnector
	daemon := daemon.New(daemonConf, db, logger.WithField("component", "daemon"))
	processes = append(processes, daemon)

	// API
	if config.API.Enabled {
		// Core API
		coreAPI := coreapi.New(&coreapi.Config{
			Database:            db,
			DefaultAuthProvider: middleware.APIKeyAuth(config.API.APIKey),
		}, logger.WithField("component", "core API"))

		// Matrix API
		matrixAPI := matrixapi.New(&matrixapi.Config{
			Database:            db,
			MatrixDB:            matrixDB,
			DefaultAuthProvider: middleware.APIKeyAuth(config.API.APIKey),
		}, logger.WithField("component", "matrix API"))

		// iCal API
		icalAPI := icalapi.New(&icalapi.Config{
			IcalDB:   icalDB,
			Database: db,
		}, logger.WithField("component", "ical API"))

		apiConfig := config.apiConfig()
		apiConfig.RouteProviders["core"] = coreAPI
		apiConfig.RouteProviders["matrix"] = matrixAPI
		apiConfig.RouteProviders["ical"] = icalAPI
		server := api.NewServer(apiConfig, logger.WithField("component", "api"))
		processes = append(processes, server)
	}

	return processes, nil
}

func assembleMatrixConfig(config *Config, icalConnector ical.Service) *matrix.Config {
	cfg := config.matrixConfig()

	cfg.DefaultMessageAction = &message.NewEventAction{}
	cfg.DefaultReplyAction = &reply.ChangeTimeAction{}

	cfg.ReplyActions = make([]matrix.ReplyAction, 0)
	cfg.MessageActions = make([]matrix.MessageAction, 0)

	cfg.MessageActions = append(cfg.MessageActions,
		&message.AddUserAction{},
		&message.EnableICalExportAction{},
	)

	cfg.BridgeServices = &matrix.BridgeServices{
		ICal: icalConnector,
	}

	return cfg
}
