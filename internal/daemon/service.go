package daemon

import (
	"sync"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

type service struct {
	config   *Config
	database database.Service
	logger   gologger.Logger
	done     chan interface{}

	daemonWG *sync.WaitGroup
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

		daemonWG: &sync.WaitGroup{},
	}
}

// Start starts the service and blocks until it either get's shut down or an un
func (service *service) Start() error {
	go service.startDailyReminderDaemon()
	service.startEventDaemon()

	service.daemonWG.Wait()

	return nil
}

func (service *service) startEventDaemon() {
	service.daemonWG.Add(1)
	eventsTicker := time.NewTicker(service.config.EventsInterval)

	for {
		select {
		case <-eventsTicker.C:
			err := service.sendOutEvents()
			if err != nil {
				service.logger.Err(err)
			}
		case <-service.done:
			service.logger.Debugf("event daemon stopped")
			service.daemonWG.Done()
			return
		}
	}
}

func (service *service) startDailyReminderDaemon() {
	service.daemonWG.Add(1)
	dailyReminderTicker := time.NewTicker(service.config.DailyReminderInterval)

	for {
		select {
		case <-dailyReminderTicker.C:
			err := service.sendOutDailyReminders()
			if err != nil {
				service.logger.Err(err)
			}
		case <-service.done:
			service.logger.Debugf("event daemon stopped")
			service.daemonWG.Done()
			return
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
