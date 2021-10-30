package log

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// Workaround for testing
func init() {
	InitLogger(true)
}

// InitLogger initializes a new logger. Make sure to call defer logger.Sync().
func InitLogger(debug bool) *zap.SugaredLogger {
	var err error
	var log *zap.Logger

	if debug {
		log, err = zap.NewDevelopment(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
	} else {
		log, err = zap.NewProduction(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
	}

	logger = log.Sugar()
	return logger
}

// Debug logs with tag debug
func Debug(msg string) {
	logger.Debug(msg)
}

// UInfo unstructured info log
func UInfo(msg string, args ...interface{}) {
	logger.Infow(msg, args...)
}

// Info logs with tag info and in blue
func Info(msg string) {
	logger.Info(msg)
}

// Warn logs with tag warn and in yellow
func Warn(msg string) {
	logger.Warn(msg)
}

// Error logs with tag error and in red
func Error(msg string) {
	logger.Error(msg)
}
