package database

import (
	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/gormlogger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type service struct {
	db     *gorm.DB
	logger gologger.Logger
	config *Config
}

//go:generate mockgen -destination=mocks/input_service.go -package=mocks . InputService

// InputService defines an interface for any arbitrary input connector.
type InputService interface {
	InputRemoved(inputType string, inputID uint, db *gorm.DB) error
}

//go:generate mockgen -destination=mocks/output_service.go -package=mocks . OutputService

// OutputService defines an interface for any arbitrary output connector.
type OutputService interface {
	OutputRemoved(outputType string, outputID uint, db *gorm.DB) error
}

type Config struct {
	LogStatements  bool
	Connection     string
	InputServices  map[string]InputService
	OutputServices map[string]OutputService
}

// NewService assembles a new database service.
func NewService(config *Config, logger gologger.Logger) (Service, error) {
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

	err = service.migrate()
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (service *service) migrate() error {
	models := []interface{}{
		Channel{},
		Input{},
		Output{},
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
	service.logger.Infof("rollbacking transaction due to error: %v", err)

	err2 := service.db.Rollback().Error
	if err2 != nil {
		service.logger.Errorf("rollbacking failed with: %v", err2)
	}

	return err
}
