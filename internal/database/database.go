package database

import (
	"time"

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

	return nil
}

// GetChannelByUserIdentifier returns the latest channel with the given user
func (d *Database) GetChannelByUserIdentifier(userID string) (*Channel, error) {
	var channel Channel
	err := d.db.First(&channel, "user_identifier = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// GetChannelByUserAndChannelIdentifier returns the latest channel with the given user and channel id
func (d *Database) GetChannelByUserAndChannelIdentifier(userID string, channelID string) (*Channel, error) {
	var channel Channel
	err := d.db.First(&channel, "user_identifier = ? AND channel_identifier = ?", userID, channelID).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// AddChannel adds a channel to the database
func (d *Database) AddChannel(userID, channelID string) (*Channel, error) {
	err := d.db.Create(&Channel{
		Created:           time.Now(),
		ChannelIdentifier: channelID,
		UserIdentifier:    userID,
	}).Error
	if err != nil {
		return nil, err
	}

	var channel Channel
	err = d.db.First(&channel, "user_identifier = ? AND channel_identifier = ?", userID, channelID).Error
	return &channel, err
}
