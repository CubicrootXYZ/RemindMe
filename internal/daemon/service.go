package daemon

import (
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

type service struct {
	Config   *Config
	Database database.Service
	Logger   gologger.Logger
	Done     chan interface{}
}

//go:generate mockgen -destination=mocks/output_service.go -package=mocks . OutputService

// OutputService defines an interface for services handling outputs.
type OutputService interface {
	SendReminder(*Event, *Output) error
	SendDailyReminder(*DailyReminder, *Output) error
}

type Config struct {
	OutputServices        map[string]OutputService // Maps OutputTypes to the services
	EventsInterval        time.Duration            // Interval in which to send out event reminders
	DailyReminderInterval time.Duration            // Interval in which to send out daily reminder
}

// NewService assembles a new service.
func NewService(config *Config, database database.Service, logger gologger.Logger) Service {
	return &service{
		Config:   config,
		Database: database,
		Logger:   logger,
		Done:     make(chan interface{}),
	}
}

// Start starts the service and blocks until it either get's shut down or an un
func (service *service) Start() error {
	eventsTicker := time.NewTicker(service.Config.EventsInterval)

	for {
		select {
		case <-eventsTicker.C:
			err := service.sendOutEvents()
			if err != nil {
				service.Logger.Err(err)
			}
		case <-service.Done:
			return nil
		}
	}
}

// Stops shuts down the service in an unblocking way.
// Start() will return once the daemon is stopped.
func (service *service) Stop() error {
	service.Logger.Debugf("stopping daemon ...")
	close(service.Done)

	return nil
}
