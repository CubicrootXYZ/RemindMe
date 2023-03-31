package database

import "gorm.io/gorm"

// TODO tests!
type service struct {
	db *gorm.DB
}

// New assembles a new iCal connector database service.
func New(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}
