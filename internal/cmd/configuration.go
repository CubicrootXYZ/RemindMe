package cmd

import (
	"flag"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/jinzhu/configor"
)

type Config struct {
	Debug    bool
	Database configDatabase
	Daemon   configDaemon
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

	return config, nil
}
