package daemon

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

func (service *service) performCleanup() error {
	opts := &database.CleanupOpts{
		OlderThan: 365 * 24 * time.Hour,
	}

	deleted, err := service.database.Cleanup(opts)
	if err != nil {
		return err
	}

	service.metricLastCleanupRun.WithLabelValues().Set(float64(time.Now().Unix()))
	service.metricItemsCleaned.WithLabelValues().Add(float64(deleted))

	return nil
}
