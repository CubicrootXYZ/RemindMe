package database

import (
	"errors"
	"strconv"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"gorm.io/gorm"
)

// Reminder is the database object for a reminder
type Reminder struct {
	gorm.Model
	RemindTime                   time.Time `gorm:"index"`
	Message                      string
	Active                       bool `gorm:"index"`
	RepeatInterval               uint64
	RepeatMax                    uint64
	Repeated                     *uint64
	ChannelID                    uint `gorm:"index"`
	Channel                      Channel
	ThirdPartyResourceID         *uint              `gorm:"index"`
	ThirdPartyResource           ThirdPartyResource `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ThirdPartyResourceIdentifier string
}

func (reminder *Reminder) GetReminderIcons() []string {
	icons := make([]string, 0)

	if reminder.RepeatInterval > 0 {
		if (reminder.Repeated != nil && *reminder.Repeated < reminder.RepeatMax) || reminder.Repeated == nil {
			icons = append(icons, "🔄")
		}
	}

	if reminder.ThirdPartyResourceID != nil {
		icons = append(icons, "📥")
	}

	return icons
}

// GET DATA

// GetReminderForChannelIDByID returns the reminder with the given ID if it relates to the given channel ID
func (d *Database) GetReminderForChannelIDByID(channelID string, reminderID int) (*Reminder, error) {
	reminder := &Reminder{}
	err := d.db.Joins("Channel").First(&reminder, "Channel.channel_identifier = ? AND reminders.id = ?", channelID, reminderID).Error

	return reminder, err
}

// GetPendingReminders returns a list with all pending reminders for the given channel
func (d *Database) GetPendingReminders(channel *Channel) ([]Reminder, error) {
	reminders := make([]Reminder, 0)

	err := d.db.Joins("Channel").Find(&reminders, "Channel.channel_identifier = ? AND active = ?", channel.ChannelIdentifier, true).Error

	return reminders, err
}

// GetMessageFromReminder returns the message with the specified message type regarding the reminder
func (d *Database) GetMessageFromReminder(reminderID uint, msgType MessageType) (*Message, error) {
	message := &Message{}
	err := d.db.Preload("Channel").Find(message, "reminder_id = ? AND type = ?", reminderID, msgType).Error
	return message, err
}

// GetPendingReminder returns all reminders that are due
func (d *Database) GetPendingReminder() ([]Reminder, error) {
	reminder := make([]Reminder, 0)
	err := d.db.Debug().Preload("Channel").Find(&reminder, "active = ? AND remind_time <= ?", 1, time.Now().UTC()).Error

	return reminder, err
}

// GetDailyReminder returns the reminders alerting in the next 24 hours
func (d *Database) GetDailyReminder(channel *Channel) (*[]Reminder, error) {
	reminder := make([]Reminder, 0)
	err := d.db.Order("remind_time asc").Joins("Channel").Find(&reminder, "channel_identifier = ? AND remind_time <= ? AND active = ?", channel.ChannelIdentifier, time.Now().Add(time.Hour*24), true).Error

	return &reminder, err
}

// UPDATE DATA

// SetReminderDone sets a reminder as inactive
func (d *Database) SetReminderDone(reminder *Reminder) (*Reminder, error) {
	if reminder.Repeated == nil {
		repeated := uint64(0)
		reminder.Repeated = &repeated
	}

	if reminder.RepeatMax > *reminder.Repeated && reminder.RepeatInterval > 0 {
		for time.Until(reminder.RemindTime) < 0 {
			reminder.RemindTime = reminder.RemindTime.Add(time.Duration(reminder.RepeatInterval) * time.Minute)
		}
		log.Debug("New remind time for reminder " + strconv.Itoa(int(reminder.ID)) + " is " + reminder.RemindTime.String())
	} else {
		reminder.Active = false
	}
	*reminder.Repeated++

	err := d.db.Save(reminder).Error

	return reminder, err
}

// UpdateReminder updates the reminder
func (d *Database) UpdateReminder(reminderID uint, remindTime time.Time, repeatInterval uint64, repeatTimes uint64) (*Reminder, error) {
	reminder := &Reminder{}
	err := d.db.First(reminder, "id = ?", reminderID).Error
	if err != nil {
		return nil, err
	}

	if time.Until(remindTime) > 0 {
		reminder.RemindTime = remindTime
	}

	if repeatInterval > 0 {
		reminder.RepeatInterval = repeatInterval
	}
	if repeatTimes > 0 {
		reminder.RepeatMax = repeatTimes
	}

	reminder.Active = true
	err = d.db.Save(reminder).Error
	return reminder, err
}

// INSERT DATA

// AddReminder adds a reminder to the database
func (d *Database) AddReminder(remindTime time.Time, message string, active bool, repeatInterval uint64, channel *Channel) (*Reminder, error) {
	reminder := Reminder{
		Message:        message,
		RemindTime:     remindTime.UTC(),
		Active:         active,
		RepeatInterval: repeatInterval,
		RepeatMax:      0,
		Channel:        *channel,
		ChannelID:      channel.ID,
	}
	reminder.Model.CreatedAt = time.Now().UTC()

	err := d.db.Create(&reminder).Error

	return &reminder, err
}

// AddOrUpdateThirdPartyResourceReminder inserts a new reminder from a third party resource or updates an existing one if already present
func (d *Database) AddOrUpdateThirdPartyResourceReminder(remindTime time.Time, message string, channelID uint, thirdPartyResourceID uint, thirdPartyResourceIdentifier string) (*Reminder, error) {
	reminder := Reminder{}
	err := d.db.First(&reminder, "channel_id = ? AND third_party_resource_id = ? AND third_party_resource_identifier = ?", channelID, thirdPartyResourceID, thirdPartyResourceIdentifier).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	reminder.Message = message
	reminder.RemindTime = remindTime
	reminder.Active = true
	reminder.ChannelID = channelID
	reminder.ThirdPartyResourceID = &thirdPartyResourceID
	reminder.ThirdPartyResourceIdentifier = thirdPartyResourceIdentifier

	err = d.db.Save(&reminder).Error
	return &reminder, err
}

// DELETE DATA

// DeleteReminder deletes the given reminder
func (d *Database) DeleteReminder(reminderID uint) (*Reminder, error) {
	deleteReminder := &Reminder{}
	err := d.db.First(&deleteReminder, "id = ?", reminderID).Error
	if err != nil {
		return nil, err
	}

	err = d.db.Delete(&deleteReminder).Statement.Error
	return deleteReminder, err
}
