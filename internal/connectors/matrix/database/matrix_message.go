package database

func (service *service) NewMessage(message *MatrixMessage) (*MatrixMessage, error) {
	err := service.db.Create(message).Error
	return message, err
}

func (service *service) GetMessageByID(messageID string) (*MatrixMessage, error) {
	var message MatrixMessage
	err := service.db.Preload("Room").Preload("User").First(&message, "matrix_messages.id = ?", messageID).Error

	return &message, err
}

func (service *service) GetLastMessage() (*MatrixMessage, error) {
	var message MatrixMessage
	err := service.db.Preload("Room").Preload("User").Order("matrix_messages.send_at DESC").First(&message).Error

	return &message, err
}
