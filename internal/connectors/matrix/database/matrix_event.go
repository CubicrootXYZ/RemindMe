package database

func (service *service) NewEvent(event *MatrixEvent) (*MatrixEvent, error) {
	err := service.db.Create(event).Error

	return event, err
}

func (service *service) GetEventByID(eventID string) (*MatrixEvent, error) {
	var event MatrixEvent
	err := service.db.Preload("Room").Preload("User").First(&event, "matrix_events.id = ?", eventID).Error

	return &event, err
}
