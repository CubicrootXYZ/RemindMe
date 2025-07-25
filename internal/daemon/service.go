package daemon

import (
	"log/slog"
	"sync"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricLastDailyReminderRun = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "remindme",
		Name:      "daemon_daily_reminder_last_run_timestamp_seconds",
		Help:      "Unix timestamp of the last run of the daily reminder daemon.",
	}, []string{})

	metricLastEventRun = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "remindme",
		Name:      "daemon_event_last_run_timestamp_seconds",
		Help:      "Unix timestamp of the last run of the event daemon.",
	}, []string{})

	metricEventsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "remindme",
		Name:      "daemon_event_events_processed_total",
		Help:      "Amount of events processed by the event daemon.",
	}, []string{})
)

type service struct {
	config   *Config
	database database.Service
	logger   *slog.Logger
	done     chan interface{}

	daemonWG *sync.WaitGroup

	metricLastDailyReminderRun *prometheus.GaugeVec
	metricLastEventRun         *prometheus.GaugeVec
	metricEventsProcessed      *prometheus.CounterVec
}

//go:generate mockgen -destination=mocks/output_service.go -package=mocks . OutputService

// OutputService defines an interface for services handling outputs.
type OutputService interface {
	ToLocalTime(time.Time, *Output) time.Time
	SendReminder(*Event, *Output) error
	SendDailyReminder(*DailyReminder, *Output) error
}

type Config struct {
	OutputServices        map[string]OutputService // Maps OutputTypes to the services
	EventsInterval        time.Duration            // Interval in which to send out event reminders
	DailyReminderInterval time.Duration            // Interval in which to send out daily reminder
}

// New assembles a new service.
func New(config *Config, database database.Service, logger *slog.Logger) Service {
	return &service{
		config:   config,
		database: database,
		logger:   logger,
		done:     make(chan interface{}),

		daemonWG: &sync.WaitGroup{},

		metricLastDailyReminderRun: metricLastDailyReminderRun,
		metricLastEventRun:         metricLastEventRun,
		metricEventsProcessed:      metricEventsProcessed,
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
			service.logger.Debug("sending out events ...")
			err := service.sendOutEvents()
			if err != nil {
				service.logger.Error("failed to send out events", "error", err)
			}
		case <-service.done:
			service.logger.Debug("event daemon stopped")
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
				service.logger.Error("failed to send out daily reminders", "error", err)
			}
		case <-service.done:
			service.logger.Debug("event daemon stopped")
			service.daemonWG.Done()
			return
		}
	}
}

// Stops shuts down the service in an unblocking way.
// Start() will return once the daemon is stopped.
func (service *service) Stop() error {
	service.logger.Debug("stopping daemon ...")
	close(service.done)

	return nil
}
