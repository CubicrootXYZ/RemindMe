package database

import (
	"errors"

	"gorm.io/gorm"
)

func (service *service) ListInputRoomsByChannel(channelID uint) ([]MatrixRoom, error) {
	var rooms []MatrixRoom
	err := service.db.Preload("Users").
		Joins("INNER JOIN inputs ON inputs.input_id = matrix_rooms.id AND inputs.input_type = ?", "matrix"). // TODO input type should be constant
		Where("inputs.channel_id = ?", channelID).
		Find(&rooms).Error

	return rooms, err
}

func (service *service) ListOutputRoomsByChannel(channelID uint) ([]MatrixRoom, error) {
	var rooms []MatrixRoom
	err := service.db.Preload("Users").
		Joins("INNER JOIN outputs ON outputs.output_id = matrix_rooms.id AND outputs.output_type = ?", "matrix"). // TODO output type should be constant
		Where("outputs.channel_id = ?", channelID).
		Find(&rooms).Error

	return rooms, err
}

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
	// Delete users associations upfront, simplest way to soft delete the room as well.
	room := &MatrixRoom{}
	room.ID = roomID
	err := service.db.Select("Users").Delete(room).Error
	if err != nil {
		return err
	}

	// Hard delete the room.
	return service.db.Unscoped().Delete(&MatrixRoom{}, "matrix_rooms.id = ?", roomID).Error
}

func (service *service) GetRoomCount() (int64, error) {
	var cnt int64
	err := service.db.Model(&MatrixRoom{}).Count(&cnt).Error

	return cnt, err
}

func (service *service) AddUserToRoom(userID string, room *MatrixRoom) (*MatrixRoom, error) {
	user, err := service.GetUserByID(userID)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return nil, err
		}
		user, err = service.NewUser(&MatrixUser{
			ID: userID,
		})
		if err != nil {
			return nil, err
		}
	}

	room.Users = append(room.Users, *user)
	room, err = service.UpdateRoom(room)
	return room, err
}
