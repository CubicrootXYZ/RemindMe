package cmd

import (
	"context"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/middleware"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	icalapi "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/api"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/reaction"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/reply"
	matrixapi "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/api"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/coreapi"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/metrics"
	"golang.org/x/sync/errgroup"
)

// Run setups the application and runs it
func Run(config *Config) error {
	logger := config.logger()

	logger.Info("starting up RemindMe", "version", config.BuildVersion)

	processes, err := setup(config, logger)
	if err != nil {
		logger.Error("startup failed", "error", err)
		return err
	}

	eg, ctx := errgroup.WithContext(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		select {
		case s := <-sigChan:
			logger.Info("shutting down", "reason", "signal", "signal", s)
		case <-ctx.Done():
			logger.Info("shutting down", "reason", "process exited")
		}

		for _, p := range processes {
			err := p.Stop()
			if err != nil {
				logger.Error("process stopped with error", "error", err)
			}
		}
	}()

	for _, p := range processes {
		eg.Go(func() error {
			return p.Start()
		})
	}

	err = eg.Wait()
	if err != nil {
		logger.Error("error group stopped with error", "error", err)
	}

	logger.Info("shut down complete, bye")

	return err
}

type process interface {
	Start() error
	Stop() error
}

func setup(config *Config, logger *slog.Logger) ([]process, error) {
	baseURL, err := url.Parse(config.API.BaseURL)
	if err != nil {
		return nil, err
	}

	processes := []process{}

	// Database
	dbConfig := config.databaseConfig()
	dbConfig.InputServices = make(map[string]database.InputService)
	dbConfig.OutputServices = make(map[string]database.OutputService)

	db, err := database.NewService(dbConfig, logger.With("component", "database"))
	if err != nil {
		logger.Error("failed to assemble database service", "error", err)
		return nil, err
	}

	// iCal connector
	icalDB, err := icaldb.New(db.GormDB())
	if err != nil {
		logger.Error("failed to assemble iCal database service", "error", err)
		return nil, err
	}

	icalConnector := ical.New(&ical.Config{
		ICalDB:          icalDB,
		Database:        db,
		BaseURL:         baseURL,
		RefreshInterval: time.Minute * time.Duration(config.ICal.RefreshInterval),
	}, logger.With("component", "ical connector"))

	dbConfig.OutputServices[ical.OutputType] = icalConnector
	dbConfig.OutputServices[ical.InputType] = icalConnector
	processes = append(processes, icalConnector)

	// Matrix connector
	matrixDB, err := matrixdb.New(db.GormDB())
	if err != nil {
		logger.Error("failed to assemble matrix database service", "error", err)
		return nil, err
	}

	matrixConnector, err := matrix.New(assembleMatrixConfig(config, icalConnector), db, matrixDB, logger.With("component", "matrix connector"))
	if err != nil {
		logger.Error("failed to assemble matrix connector service", "error", err)
		return nil, err
	}

	processes = append(processes, matrixConnector)

	dbConfig.InputServices[matrix.InputType] = matrixConnector
	dbConfig.OutputServices[matrix.OutputType] = matrixConnector

	// Daemon
	daemonConf := config.daemonConfig()
	daemonConf.OutputServices = make(map[string]daemon.OutputService)
	daemonConf.OutputServices[matrix.OutputType] = matrixConnector
	daemonConf.OutputServices[ical.OutputType] = icalConnector
	daemon := daemon.New(daemonConf, db, logger.With("component", "daemon"))
	processes = append(processes, daemon)

	// API
	if config.API.Enabled {
		// Core API
		coreAPI := coreapi.New(&coreapi.Config{
			Database:            db,
			DefaultAuthProvider: middleware.APIKeyAuth(config.API.APIKey),
		}, logger.With("component", "core API"))

		// Matrix API
		matrixAPI := matrixapi.New(&matrixapi.Config{
			Database:            db,
			MatrixDB:            matrixDB,
			DefaultAuthProvider: middleware.APIKeyAuth(config.API.APIKey),
		}, logger.With("component", "matrix API"))

		// iCal API
		icalAPI := icalapi.New(&icalapi.Config{
			IcalDB:   icalDB,
			Database: db,
		}, logger.With("component", "ical API"))

		apiConfig := config.apiConfig()
		apiConfig.RouteProviders["core"] = coreAPI
		apiConfig.RouteProviders["matrix"] = matrixAPI
		apiConfig.RouteProviders["ical"] = icalAPI
		server := api.NewServer(apiConfig, logger.With("component", "api"))
		processes = append(processes, server)
	}

	// Metrics
	if config.Metrics.Enabled {
		metricsService, err := metrics.New(config.metricsConfig())
		if err != nil {
			return nil, err
		}

		processes = append(processes, metricsService)
	}

	return processes, nil
}

func assembleMatrixConfig(config *Config, icalConnector ical.Service) *matrix.Config {
	cfg := config.matrixConfig()

	cfg.DefaultMessageAction = &message.NewEventAction{}
	cfg.DefaultReplyAction = &reply.ChangeTimeAction{}

	cfg.ReplyActions = make([]matrix.ReplyAction, 0)
	cfg.MessageActions = make([]matrix.MessageAction, 0)
	cfg.ReactionActions = make([]matrix.ReactionAction, 0)

	cfg.ReplyActions = append(cfg.ReplyActions,
		&reply.DeleteEventAction{},
		&reply.MakeRecurringAction{},
	)

	cfg.MessageActions = append(cfg.MessageActions,
		&message.AddUserAction{},
		&message.EnableICalExportAction{},
		&message.ChangeTimezoneAction{},
		&message.RegenICalTokenAction{},
		&message.DeleteEventAction{},
		&message.SetDailyReminderAction{},
		&message.ListEventsAction{},
		&message.RegenICalTokenAction{},
		&message.ChangeEventAction{},
		&message.ListCommandsAction{},
	)

	cfg.ReactionActions = append(cfg.ReactionActions,
		&reaction.DeleteEventAction{},
		&reaction.AddTimeAction{},
		&reaction.MarkDoneAction{},
		&reaction.RescheduleRepeatingAction{},
	)

	cfg.BridgeServices = &matrix.BridgeServices{
		ICal: icalConnector,
	}

	return cfg
}
