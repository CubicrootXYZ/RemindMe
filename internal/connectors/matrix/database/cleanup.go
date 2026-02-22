package database

import "time"

func (service *service) Cleanup() error {
	return service.db.Delete(&MatrixMessage{}, "send_at < ? AND event_id IS NULL", time.Now().Add(-180*24*time.Hour)).Error
}
