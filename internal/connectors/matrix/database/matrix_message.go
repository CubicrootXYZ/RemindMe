package database

import (
	"errors"

	"gorm.io/gorm"
)

func (service *service) NewMessage(message *MatrixMessage) (*MatrixMessage, error) {
	err := service.db.Create(message).Error
	return message, err
}

func (service *service) GetMessageByID(messageID string) (*MatrixMessage, error) {
	var message MatrixMessage
	err := service.db.Preload("Room").Preload("User").First(&message, "matrix_messages.id = ?", messageID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &message, err
}

func (service *service) GetLastMessage() (*MatrixMessage, error) {
	var message MatrixMessage
	err := service.db.Preload("Room").Preload("User").Order("matrix_messages.send_at DESC").First(&message).Error

	return &message, err
}

func (service *service) DeleteAllMessagesFromRoom(roomID uint) error {
	return service.db.Unscoped().Delete(&MatrixMessage{}, "matrix_messages.room_id = ?", roomID).Error
}
