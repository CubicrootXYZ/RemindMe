package ical

import icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"

// The in- and output type provided by this package
const (
	InputType  = "ical"
	OutputType = "ical"
)

// Service provides and interface for the ical connector.
// The connector is suitable for in- and output.
type Service interface {
	Start() error
	Stop() error

	InputRemoved(inputType string, inputID uint) error
	OutputRemoved(outputType string, outputID uint) error

	NewOutput(channelID uint) (*icaldb.IcalOutput, error)
}
