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
	DailyReminder     *uint // minutes from midnight when to send the daily reminder. Null to deactivate.
}

// Timezone returns the timezone of the channel
func (c *Channel) Timezone() *time.Location {
	if c.TimeZone == "" {
		return time.UTC
	}
	loc, err := time.LoadLocation(c.TimeZone)
	if err != nil {
		return time.UTC
	}

	return loc
}

// GET DATA

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

// GetChannelList returns all known channels
func (d *Database) GetChannelList() ([]Channel, error) {
	channel := make([]Channel, 0)

	err := d.db.Model(&channel).Find(&channel).Error

	return channel, err
}

// UPDATE DATA

// UpdateChannel updates the given channel with the given information
func (d *Database) UpdateChannel(channelID uint, timeZone string, dailyReminder *uint) (*Channel, error) {
	channel := &Channel{}
	err := d.db.First(channel, "id = ?", channelID).Error
	if err != nil {
		return nil, err
	}

	channel.TimeZone = timeZone
	if dailyReminder != nil {
		channel.DailyReminder = dailyReminder
	}

	err = d.db.Save(channel).Error
	return channel, err
}

// INSERT DATA

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
