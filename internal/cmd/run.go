package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/reply"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
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
	dbConfig := config.databaseConfig()
	db, err := database.NewService(dbConfig, logger.WithField("component", "database"))
	if err != nil {
		logger.Err(err)
		return nil, err
	}

	// TODO set message and reply actions
	matrixDB, err := matrixdb.New(db.GormDB())
	if err != nil {
		logger.Err(err)
		return nil, err
	}

	matrixConnector, err := matrix.New(assembleMatrixConfig(config), db, matrixDB, logger.WithField("component", "matrix connector"))
	if err != nil {
		log.Err(err)
		return nil, err
	}

	// TODO move services to matrixDB or own interface to remove circular dependency db => matrixCon => db
	dbConfig.InputServices = make(map[string]database.InputService)
	dbConfig.InputServices[matrix.InputType] = matrixConnector
	dbConfig.InputServices[matrix.OutputType] = matrixConnector

	daemonConf := config.daemonConfig()
	daemonConf.OutputServices = make(map[string]daemon.OutputService)
	// TODO daemonConf.OutputServices[matrix.OutputType] = matrixConnector
	daemon := daemon.New(daemonConf, db, logger.WithField("component", "daemon"))

	return []process{daemon, matrixConnector}, nil
}

func assembleMatrixConfig(config *Config) *matrix.Config {
	cfg := config.matrixConfig()

	cfg.DefaultMessageAction = &message.NewEventAction{}
	cfg.DefaultReplyAction = &reply.ChangeTimeAction{}

	cfg.ReplyActions = make([]matrix.ReplyAction, 0)
	cfg.MessageActions = make([]matrix.MessageAction, 0)

	cfg.MessageActions = append(cfg.MessageActions,
		&message.AddUserAction{},
	)

	return cfg
}
