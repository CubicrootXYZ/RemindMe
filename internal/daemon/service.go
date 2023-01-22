package daemon

import (
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

type service struct {
	config   *Config
	database database.Service
	logger   gologger.Logger
	done     chan interface{}
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

// New assembles a new service.
func New(config *Config, database database.Service, logger gologger.Logger) Service {
	return &service{
		config:   config,
		database: database,
		logger:   logger,
		done:     make(chan interface{}),
	}
}

// Start starts the service and blocks until it either get's shut down or an un
func (service *service) Start() error {
	eventsTicker := time.NewTicker(service.config.EventsInterval)
	dailyReminderTicker := time.NewTicker(service.config.DailyReminderInterval)

	for {
		select {
		case <-eventsTicker.C:
			err := service.sendOutEvents()
			if err != nil {
				service.logger.Err(err)
			}
		case <-dailyReminderTicker.C:
			err := service.sendOutDailyReminders()
			if err != nil {
				service.logger.Err(err)
			}
		case <-service.done:
			service.logger.Debugf("daemon stopped")
			return nil
		}
	}
}

// Stops shuts down the service in an unblocking way.
// Start() will return once the daemon is stopped.
func (service *service) Stop() error {
	service.logger.Debugf("stopping daemon ...")
	close(service.done)

	return nil
}
