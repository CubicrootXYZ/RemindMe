package database

import "gorm.io/gorm"

// TODO tests!
type service struct {
	db *gorm.DB
}

// New assembles a new iCal connector database service.
func New(db *gorm.DB) (Service, error) {
	s := &service{
		db: db,
	}

	err := s.migrate()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (service *service) migrate() error {
	err := service.db.AutoMigrate(&IcalInput{})
	if err != nil {
		return err
	}
	err = service.db.AutoMigrate(&IcalOutput{})
	if err != nil {
		return err
	}
	return nil
}
