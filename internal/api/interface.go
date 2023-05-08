package api

// Server defines the API webserver interface.
type Server interface {
	Start() error
	Stop() error
}
