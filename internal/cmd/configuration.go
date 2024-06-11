package cmd

import (
	"errors"
	"flag"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/jinzhu/configor"
)

type Config struct {
	Debug    bool
	Database configDatabase
	Daemon   configDaemon
	Matrix   configMatrix
	ICal     configICal
	API      configAPI

	BuildVersion string
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

func (config *Config) databaseConfig() *database.Config {
	return &database.Config{
		Connection: config.Database.Connection,
	}
}

func (config *Config) loggerConfig() gologger.LogLevel {
	if config.Debug {
		return gologger.LogLevelDebug
	}

	return gologger.LogLevelInfo
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
		EnableE2EE:    config.Matrix.Bot.E2EE,
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

	return config, nil
}
