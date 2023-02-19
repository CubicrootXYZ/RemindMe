package matrix

import (
	"errors"
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

// Service provides and interface for the matrix connectore.
// The connector is suitable for in- and output.
type Service interface {
	Start() error
	Stop() error

	InputRemoved(inputType string, inputID uint) error
	OutputRemoved(outputType string, outputID uint) error
}
