package database

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/random"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"gorm.io/gorm"
)

// Channel holds data about a messaging channel
type Channel struct {
	gorm.Model
	Created           time.Time
	ChannelIdentifier string `gorm:"index;size:500"`
	UserIdentifier    string `gorm:"index;size:500"`
	TimeZone          string
	DailyReminder     *uint  // minutes from midnight when to send the daily reminder. Null to deactivate.
	CalendarSecret    string `gorm:"index"`
	Role              *roles.Role
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

// GetChannel returns the channel
func (d *Database) GetChannel(id uint) (*Channel, error) {
	var channel Channel
	err := d.db.First(&channel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
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

// GetChannelsByChannelIdentifier returns all channels with the given channel identifier
func (d *Database) GetChannelsByChannelIdentifier(channelID string) ([]Channel, error) {
	channels := make([]Channel, 0)
	err := d.db.Find(&channels, "channel_identifier = ?", channelID).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
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
func (d *Database) UpdateChannel(channelID uint, timeZone string, dailyReminder *uint, role *roles.Role) (*Channel, error) {
	channel := &Channel{}
	err := d.db.First(channel, "id = ?", channelID).Error
	if err != nil {
		return nil, err
	}

	channel.TimeZone = timeZone
	channel.DailyReminder = dailyReminder
	channel.Role = role

	err = d.db.Save(channel).Error
	return channel, err
}

// GenerateNewCalendarSecret generates and sets a new calendar secret
func (d *Database) GenerateNewCalendarSecret(channel *Channel) error {
	channel.CalendarSecret = random.String(30)
	return d.db.Save(&channel).Error
}

// INSERT DATA

// AddChannel adds a channel to the database
func (d *Database) AddChannel(userID, channelID string, role roles.Role) (*Channel, error) {
	err := d.db.Create(&Channel{
		Created:           time.Now(),
		ChannelIdentifier: channelID,
		UserIdentifier:    userID,
		Role:              &role,
		CalendarSecret:    random.String(30),
	}).Error
	if err != nil {
		return nil, err
	}

	var channel Channel
	err = d.db.First(&channel, "user_identifier = ? AND channel_identifier = ?", userID, channelID).Error
	return &channel, err
}

// DELETE DATA

// CleanAdminChannels removes all admin channels except the ones given in keep
func (d *Database) CleanAdminChannels(keep []*Channel) error {
	channels, err := d.GetChannelList()
	if err != nil {
		return err
	}

	for _, channel := range channels {
		remove := true
		for _, channelKeep := range keep {
			if channel.ID == channelKeep.ID && channel.ChannelIdentifier == channelKeep.ChannelIdentifier && channel.UserIdentifier == channelKeep.UserIdentifier {
				remove = false
				break
			}
		}

		if channel.Role != nil && *channel.Role != roles.RoleAdmin {
			remove = false
		}

		if remove {
			log.Info(fmt.Sprintf("Removing channel %d", channel.ID))
			err = d.db.Model(&Reminder{}).Where("channel_id = ?", channel.ID).Updates(map[string]interface{}{"active": 0}).Error
			if err != nil {
				return err
			}
			err = d.db.Delete(&channel).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteChannel deletes the given channel
func (d *Database) DeleteChannel(channel *Channel) error {
	return d.db.Delete(channel).Error
}
