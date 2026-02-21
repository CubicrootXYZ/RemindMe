package database

import "time"

type CleanupOpts struct {
	OlderThan time.Duration
}

func (service *service) Cleanup(opts *CleanupOpts) (int64, error) {
	result := service.db.Unscoped().Where("time < ? AND (repeat_until IS NULL OR repeat_until < ?)", time.Now().Add(-opts.OlderThan), time.Now().Add(-opts.OlderThan)).Delete(&Event{})

	return result.RowsAffected, result.Error
}
