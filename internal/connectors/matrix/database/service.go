package database

import "gorm.io/gorm"

type service struct {
	db *gorm.DB
}

func New(db *gorm.DB) (Service, error) {
	service := service{
		db: db,
	}

	err := service.migrate()
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (service *service) migrate() error {
	models := []interface{}{
		MatrixRoom{},
		MatrixUser{},
	}

	for _, m := range models {
		err := service.db.AutoMigrate(m)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *service) UpdateRoom(room *MatrixRoom) (*MatrixRoom, error) {
	err := service.db.Save(room).Error

	return room, err
}

func (service *service) GetUserByID(userID string) (*MatrixUser, error) {
	var user MatrixUser
	err := service.db.Preload("Rooms").First(&user, "id = ?", userID).Error

	return &user, err
}
