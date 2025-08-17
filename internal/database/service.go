package database

import (
	"log/slog"

	"github.com/CubicrootXYZ/gormlogger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type service struct {
	db     *gorm.DB
	logger *slog.Logger
	config *Config
}

//go:generate mockgen -destination=mocks/input_service.go -package=mocks . InputService

// InputService defines an interface for any arbitrary input connector.
type InputService interface {
	InputRemoved(inputType string, inputID uint) error
}

//go:generate mockgen -destination=mocks/output_service.go -package=mocks . OutputService

// OutputService defines an interface for any arbitrary output connector.
type OutputService interface {
	OutputRemoved(outputType string, outputID uint) error
}

type Config struct {
	LogStatements  bool
	Connection     string
	InputServices  map[string]InputService
	OutputServices map[string]OutputService
}

// NewService assembles a new database service.
func NewService(config *Config, logger *slog.Logger) (Service, error) {
	logger.Debug("setting up database")

	if config == nil {
		return nil, ErrInvalidConfig
	}

	gormLogger := gormlogger.NewLogger(config.LogStatements)

	db, err := gorm.Open(mysql.Open(config.Connection+"?parseTime=True"), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	service := &service{
		config: config,
		db:     db,
		logger: logger,
	}

	logger.Debug("migrating database")

	err = service.migrate()
	if err != nil {
		return nil, err
	}

	logger.Debug("database setup finished")

	return service, nil
}

func (service *service) GormDB() *gorm.DB {
	return service.db
}

func (service *service) migrate() error {
	models := []interface{}{
		Channel{},
		Input{},
		Output{},
		Event{},
	}

	for i := range models {
		err := service.db.AutoMigrate(&models[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *service) newSession() *service {
	tx := service.db.Session(&gorm.Session{
		SkipDefaultTransaction: true,
	}).Begin()

	newService := *service
	newService.db = tx

	return &newService
}

func (service *service) commit() error {
	return service.db.Commit().Error
}

func (service *service) rollbackWithError(err error) error {
	service.logger.Info("rollback transaction", "error", err)

	err2 := service.db.Rollback().Error
	if err2 != nil {
		service.logger.Error("rollback failed", "error", err2, "original_error", err)
	}

	return err
}
