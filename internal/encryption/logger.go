package encryption

import (
	"go.uber.org/zap"
)

type cryptoLogger struct {
	log *zap.SugaredLogger
}

func newCryptoLogger(debug bool) (*cryptoLogger, error) {
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

	logger := log.Sugar()

	return &cryptoLogger{
		log: logger,
	}, nil
}

func (c cryptoLogger) Error(message string, args ...interface{}) {
	c.log.Errorf(message, args...)
}

func (c cryptoLogger) Warn(message string, args ...interface{}) {
	c.log.Warnf(message, args...)
}

func (c cryptoLogger) Debug(message string, args ...interface{}) {
	c.log.Debugf(message, args...)
}

func (c cryptoLogger) Trace(message string, args ...interface{}) {
	c.log.Debugf(message, args...)
}
