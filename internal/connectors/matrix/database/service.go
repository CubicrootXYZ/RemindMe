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
		MatrixMessage{},
		MatrixEvent{},
	}

	for _, m := range models {
		err := service.db.AutoMigrate(m)
		if err != nil {
			return err
		}
	}

	return nil
}
