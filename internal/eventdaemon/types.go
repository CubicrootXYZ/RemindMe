package eventdaemon

// Syncer is responsible for receiving messages from a messenger
type Syncer interface {
	Start(daemon *Daemon) error
	Stop()
}
