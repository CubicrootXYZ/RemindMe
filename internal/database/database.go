package database

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database holds all information for connecting to the database
type Database struct {
	config configuration.Database
	db     gorm.DB
}

// Create creates a database object
func Create(config configuration.Database) (*Database, error) {
	db := Database{
		config: config,
	}

	err := db.initialize()

	return &db, err
}

func (d *Database) initialize() error {
	db, err := gorm.Open(mysql.Open(d.config.Connection+"?parseTime=True"), &gorm.Config{})
	if err != nil {
		return err
	}

	d.db = *db

	err = d.db.AutoMigrate(&Reminder{})
	if err != nil {
		return err
	}

	err = d.db.AutoMigrate(&Channel{})
	if err != nil {
		return err
	}

	err = d.db.AutoMigrate(&Message{})
	if err != nil {
		return err
	}

	err = d.db.AutoMigrate(&Event{})
	if err != nil {
		return err
	}

	return nil
}
