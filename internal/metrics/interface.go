package metrics

// Service defines an interface for a metrics service.
type Service interface {
	Start() error
	Stop() error
}
