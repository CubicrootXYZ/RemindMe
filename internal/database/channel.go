package database

import (
	"time"

	"gorm.io/gorm"
)

// Channel holds data about a messaging channel
type Channel struct {
	gorm.Model
	Created           time.Time
	ChannelIdentifier string
	UserIdentifier    string
	TimeZone          string
}

// UpdateChannel updates the given channel with the given information
func (d *Database) UpdateChannel(channelID uint, timeZone string) (*Channel, error) {
	channel := &Channel{}
	err := d.db.First(channel, "id = ?", channelID).Error
	if err != nil {
		return nil, err
	}

	channel.TimeZone = timeZone
	err = d.db.Save(channel).Error
	return channel, err
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
