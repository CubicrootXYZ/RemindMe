package matrix

import (
	"errors"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
)

// The in- and output type provided by this package
const (
	InputType  = "matrix"
	OutputType = "matrix"
)

// Errors exposed by the package.
var (
	ErrUnknowEvent = errors.New("unknown event")
)

// Service provides and interface for the matrix connector.
// The connector is suitable for in- and output.
type Service interface {
	Start() error
	Stop() error

	InputRemoved(inputType string, inputID uint) error
	OutputRemoved(outputType string, outputID uint) error

	SendDailyReminder(*daemon.DailyReminder, *daemon.Output) error
	SendReminder(*daemon.Event, *daemon.Output) error

	ToLocalTime(time.Time, *daemon.Output) time.Time
}
