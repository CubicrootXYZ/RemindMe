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
	err := service.db.Preload("Event").Preload("Room").Preload("User").First(&message, "matrix_messages.id = ?", messageID).Error
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
	err := service.db.Preload("Event").Preload("Room").Preload("User").Order("matrix_messages.send_at DESC").First(&message).Error

	return &message, err
}

func (service *service) GetEventMessageByOutputAndEvent(eventID uint, outputID uint, _ string) (*MatrixMessage, error) {
	var room MatrixRoom
	err := service.db.First(&room, "matrix_rooms.id = ?", outputID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var message MatrixMessage
	err = service.db.Preload("Event").Preload("Room").Preload("User").
		First(
			&message,
			"matrix_messages.event_id = ? AND matrix_messages.type = ? AND matrix_messages.incoming = ? AND matrix_messages.room_id = ?",
			eventID,
			MessageTypeNewEvent,
			true,
			room.ID,
		).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &message, nil
}

func (service *service) DeleteAllMessagesFromRoom(roomID uint) error {
	return service.db.Unscoped().Delete(&MatrixMessage{}, "matrix_messages.room_id = ?", roomID).Error
}
