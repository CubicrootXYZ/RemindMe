package cmd

import (
	"errors"
	"flag"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/CubicrootXYZ/configor"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/metrics"
	"github.com/lmittmann/tint"
)

type Config struct {
	Database configDatabase
	Daemon   configDaemon
	Matrix   configMatrix
	ICal     configICal
	API      configAPI
	Logger   configLogger
	Metrics  configMetrics

	BuildVersion string
}

type configLogger struct {
	Format string `default:"text"`
	Debug  bool
}

type configDatabase struct {
	Connection    string `required:"true"`
	LogStatements bool
}

type configDaemon struct {
	Intervals struct {
		Events         uint `default:"30"`
		DailyReminders uint `default:"600"`
	}
}

type configMatrix struct {
	Bot struct {
		Username   string
		Password   string
		Homeserver string
		DeviceID   string
		E2EE       bool
		DeviceKey  string
	}
	AllowInvites  bool
	RoomLimit     uint
	UserWhitelist []string
}

type configICal struct {
	RefreshInterval uint `default:"60"`
}

type configAPI struct {
	Enabled bool
	Address string `default:"0.0.0.0:8080"`
	APIKey  string
	BaseURL string
}

type configMetrics struct {
	Enabled bool
	Address string `default:"0.0.0.0:9092"`
}

func (config *Config) databaseConfig() *database.Config {
	return &database.Config{
		Connection: config.Database.Connection,
	}
}

func (config *Config) logger() *slog.Logger {
	logLevel := slog.LevelInfo
	if config.Logger.Debug {
		logLevel = slog.LevelDebug
	}

	var handler slog.Handler
	switch config.Logger.Format {
	case "text":
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      logLevel,
			TimeFormat: time.RFC3339Nano,
		})
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	return slog.New(handler)
}

func (config *Config) daemonConfig() *daemon.Config {
	return &daemon.Config{
		EventsInterval:        time.Second * time.Duration(config.Daemon.Intervals.Events),
		DailyReminderInterval: time.Second * time.Duration(config.Daemon.Intervals.DailyReminders),
	}
}

func (config *Config) matrixConfig() *matrix.Config {
	return &matrix.Config{
		Username:      config.Matrix.Bot.Username,
		Password:      config.Matrix.Bot.Password,
		Homeserver:    config.Matrix.Bot.Homeserver,
		DeviceID:      config.Matrix.Bot.DeviceID,
		DeviceKey:     config.Matrix.Bot.DeviceKey,
		AllowInvites:  config.Matrix.AllowInvites,
		RoomLimit:     config.Matrix.RoomLimit,
		UserWhitelist: config.Matrix.UserWhitelist,
	}
}

func (config *Config) apiConfig() *api.Config {
	return &api.Config{
		Address:        config.API.Address,
		RouteProviders: make(map[string]api.RouteProvider),
	}
}

func (config *Config) metricsConfig() *metrics.Config {
	return &metrics.Config{
		Address: config.Metrics.Address,
	}
}

func LoadConfiguration() (*Config, error) {
	fileName := flag.String("config", "config.yml", "Configuration file to load")
	flag.Parse()

	config := &Config{}

	err := configor.New(&configor.Config{
		Environment:          "production",
		ENVPrefix:            "REMINDME",
		ErrorOnUnmatchedKeys: false,
	}).Load(config, *fileName)
	if err != nil {
		return nil, err
	}

	if config.API.Enabled && len(config.API.APIKey) < 10 {
		return nil, errors.New("API key needs to be at least 10 characters")
	}

	if !slices.Contains([]string{"text", "json"}, config.Logger.Format) {
		return nil, errors.New("logger format must be one of: text, json")
	}

	return config, nil
}
