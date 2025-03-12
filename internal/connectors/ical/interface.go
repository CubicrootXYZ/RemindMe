package ical

import (
	"errors"
	"time"

	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
)

// The in- and output type provided by this package
const (
	InputType  = "ical"
	OutputType = "ical"
)

// List of commonly used errors in this package.
var (
	ErrNotFound = errors.New("not found")
)

//go:generate mockgen -destination=service_mock.go -package=ical . Service

// Service provides and interface for the ical connector.
// The connector is suitable for in- and output.
type Service interface {
	Start() error
	Stop() error

	InputRemoved(inputType string, inputID uint) error
	OutputRemoved(outputType string, outputID uint) error

	NewOutput(channelID uint) (*icaldb.IcalOutput, string, error)                 // Returns the calendar URL.
	GetOutput(outputID uint, regenToken bool) (*icaldb.IcalOutput, string, error) // Returns the calendar URL.

	SendReminder(*daemon.Event, *daemon.Output) error
	SendDailyReminder(*daemon.DailyReminder, *daemon.Output) error

	ToLocalTime(time.Time, *daemon.Output) time.Time
}
