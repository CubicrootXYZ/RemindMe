package database

import (
	"errors"

	"gorm.io/gorm"
)

func (service *service) GetRoomByID(roomID string) (*MatrixRoom, error) {
	var room MatrixRoom
	err := service.db.Preload("Users").First(&room, "room_id = ?", roomID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &room, err
}

func (service *service) NewRoom(room *MatrixRoom) (*MatrixRoom, error) {
	err := service.db.Create(room).Error

	return room, err
}

func (service *service) UpdateRoom(room *MatrixRoom) (*MatrixRoom, error) {
	err := service.db.Save(room).Error

	return room, err
}
