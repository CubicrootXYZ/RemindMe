package database

import (
	"errors"

	"gorm.io/gorm"
)

func (service *service) NewEvent(event *MatrixEvent) (*MatrixEvent, error) {
	err := service.db.Create(event).Error

	return event, err
}

func (service *service) GetEventByID(eventID string) (*MatrixEvent, error) {
	var event MatrixEvent
	err := service.db.Preload("Room").Preload("User").First(&event, "matrix_events.id = ?", eventID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &event, err
}

func (service *service) DeleteAllEventsFromRoom(roomID uint) error {
	return service.db.Unscoped().Delete(&MatrixEvent{}, "matrix_events.room_id = ?", roomID).Error
}
