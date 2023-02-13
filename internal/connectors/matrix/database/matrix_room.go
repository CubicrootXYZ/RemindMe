package database

import (
	"errors"

	"gorm.io/gorm"
)

func (service *service) GetRoomByID(id uint) (*MatrixRoom, error) {
	var room MatrixRoom
	err := service.db.Preload("Users").First(&room, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &room, err
}

func (service *service) GetRoomByRoomID(roomID string) (*MatrixRoom, error) {
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

func (service *service) DeleteRoom(roomID uint) error {
	return service.db.Unscoped().Delete(&MatrixRoom{}, "matrix_rooms.id = ?", roomID).Error
}

func (service *service) GetRoomCount() (int64, error) {
	var cnt int64
	err := service.db.Model(&MatrixRoom{}).Count(&cnt).Error

	return cnt, err
}
