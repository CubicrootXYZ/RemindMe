package matrix

// OutputService defines a service that acts as an output.
type OutputService interface {
}

// InputService defines a service that acts as an input.
type InputService interface {
	Start() error
	Stop() error
}

// The in- and output type provided by this package
const (
	InputType  = "matrix"
	OutputType = "matrix"
)
