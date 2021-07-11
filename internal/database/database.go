package database

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/event"
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

// AddReminder adds a reminder to the database
func (d *Database) AddReminder(remindTime time.Time, message string, active bool, repeatInterval uint64, channel *Channel) (*Reminder, error) {
	reminder := Reminder{
		Message:        message,
		RemindTime:     remindTime,
		Active:         true,
		RepeatInterval: repeatInterval,
		RepeatMax:      0,
		Channel:        *channel,
		ChannelID:      channel.ID,
	}
	reminder.Model.CreatedAt = time.Now().UTC()

	err := d.db.Create(&reminder).Error

	return &reminder, err
}

// AddMessage adds a message to the database
func (d *Database) AddMessage(id string, timestamp int64, content *event.MessageEventContent, reminder *Reminder, msgType MessageType, channel *Channel) (*Message, error) {
	relatesTo := ""
	if content.RelatesTo != nil {
		relatesTo = content.RelatesTo.EventID.String()
	}
	message := Message{
		Body:               content.Body,
		BodyHTML:           content.FormattedBody,
		ReminderID:         reminder.ID,
		Reminder:           *reminder,
		ResponseToMessage:  relatesTo,
		Type:               msgType,
		ChannelID:          channel.ID,
		Channel:            *channel,
		Timestamp:          timestamp,
		ExternalIdentifier: id,
	}
	message.Model.CreatedAt = time.Now().UTC()

	err := d.db.Create(&message).Error

	return &message, err
}
