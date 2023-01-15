package eventdaemon

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

// Daemon is the event daemon collecting events from a messenger
type Daemon struct {
	Database types.Database
	syncer   Syncer
	done     chan interface{}
}

// Command defines commands the daemon can handle
type Command int16

const (
	// CommandStop stops the daemon
	CommandStop = Command(1)
)

// Create returns a new event daemon
func Create(database types.Database, syncer Syncer) *Daemon {
	return &Daemon{
		Database: database,
		syncer:   syncer,
		done:     make(chan interface{}),
	}
}

// Start starts the daemon
func (d *Daemon) Start() error {
	for {
		log.Info("starting matrix syncer")
		err := d.syncer.Start(d)
		if err != nil {
			log.Error("matrix syncer returned error: " + err.Error())
		}

		select {
		case <-d.done:
			log.Info("event daemon stopped")
			return nil
		default:
		}

		log.Warn("matrix syncer stopped, starting again in 2 minutes")
		d.syncer.Stop()
		time.Sleep(time.Minute * 2)
	}
}

// Stop stops the daemon
func (d *Daemon) Stop() {
	log.Debug("stopping event daemon ...")
	d.syncer.Stop()
	close(d.done)
}
