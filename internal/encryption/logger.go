package encryption

import (
	"fmt"
	"log"
)

type cryptoLogger struct {
	prefix string
}

func (c cryptoLogger) Error(message string, args ...interface{}) {
	log.Printf(fmt.Sprintf("[%s/Error] %s", c.prefix, message), args...)
}

func (c cryptoLogger) Warn(message string, args ...interface{}) {
	log.Printf(fmt.Sprintf("[%s/Warn] %s", c.prefix, message), args...)
}

func (c cryptoLogger) Debug(message string, args ...interface{}) {
	log.Printf(fmt.Sprintf("[%s/Debug] %s", c.prefix, message), args...)
}

func (c cryptoLogger) Trace(message string, args ...interface{}) {
	log.Printf(fmt.Sprintf("[%s/Trace] %s", c.prefix, message), args...)
}
