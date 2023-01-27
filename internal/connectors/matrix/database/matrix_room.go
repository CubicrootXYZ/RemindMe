package database

func (service *service) GetRoomByID(roomID string) (*MatrixRoom, error) {
	var room MatrixRoom
	err := service.db.Preload("Users").First(&room, "room_id = ?", roomID).Error

	return &room, err
}

func (service *service) NewRoom(room *MatrixRoom) (*MatrixRoom, error) {
	err := service.db.Create(room).Error

	return room, err
}
