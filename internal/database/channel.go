package database

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/random"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/id"
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
	LastCryptoEvent   string `gorm:"type:text"`
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

// GetChannelsByUserIdentifier returns all channels with the given user
func (d *Database) GetChannelsByUserIdentifier(userID string) ([]Channel, error) {
	channels := make([]Channel, 0)
	err := d.db.Find(&channels, "user_identifier = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
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

// ChannelCount returns the amount of active channels
func (d *Database) ChannelCount() (int64, error) {
	var count int64
	err := d.db.Model(&Channel{}).Count(&count).Error
	return count, err
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

// ChannelSaveChanges saves the changes in the given channel
func (d *Database) ChannelSaveChanges(channel *Channel) error {
	return d.db.Save(channel).Error
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

	for i := range channels {
		remove := true
		for _, channelKeep := range keep {
			if channels[i].ID == channelKeep.ID && channels[i].ChannelIdentifier == channelKeep.ChannelIdentifier && channels[i].UserIdentifier == channelKeep.UserIdentifier {
				remove = false
				break
			}
		}

		if channels[i].Role != nil && *channels[i].Role != roles.RoleAdmin {
			remove = false
		}

		if remove {
			log.Info(fmt.Sprintf("Removing channel %d", channels[i].ID))
			err = d.DeleteChannel(&channels[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteChannel deletes the given channel
func (d *Database) DeleteChannel(channel *Channel) error {
	channels, err := d.GetChannelsByChannelIdentifier(channel.ChannelIdentifier)
	if err != nil {
		return err
	}

	if d.matrixClient != nil && len(channels) == 1 {
		_, err := d.matrixClient.LeaveRoom(id.RoomID(channel.ChannelIdentifier))
		if err != nil {
			log.Warn("Failed to leave room with: " + err.Error())
		}
	}

	err = d.db.Unscoped().Delete(&Message{}, "channel_id = ?", channel.ID).Error
	if err != nil {
		return err
	}

	err = d.db.Unscoped().Delete(&Reminder{}, "channel_id = ?", channel.ID).Error
	if err != nil {
		return err
	}

	err = d.db.Model(&Event{}).Where("channel_id = ?", channel.ID).Updates(map[string]interface{}{"channel_id": nil}).Error
	if err != nil {
		return err
	}

	return d.db.Unscoped().Delete(channel).Error
}

// DeleteChannelsFromUser removed all channels from the given matrix user
func (d *Database) DeleteChannelsFromUser(userID string) error {
	channels := make([]Channel, 0)
	err := d.db.Find(&channels, "user_identifier = ? AND role != 'admin'", userID).Error
	if err != nil {
		return err
	}

	for i := range channels {
		err := d.DeleteChannel(&channels[i])
		if err != nil {
			return err
		}
	}

	return nil
}
