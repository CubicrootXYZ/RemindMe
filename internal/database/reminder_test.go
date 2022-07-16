package database

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReminder_GetReminderForChannelIDByIDOnSuccess(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, reminder := range testReminders() {
		mock.ExpectQuery("SELECT (.*) FROM `reminders`").
			WithArgs(
				"!abcdefghijklmop",
				reminder.ID,
			).
			WillReturnRows(rowsForReminders([]*Reminder{reminder}))

		c := &Channel{}
		c.ID = reminder.ChannelID
		retReminder, err := db.GetReminderForChannelIDByID("!abcdefghijklmop", int(reminder.ID))
		require.NoError(t, err)

		require.False(t, retReminder == nil)
		assert.Equal(reminder.RemindTime, retReminder.RemindTime)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestReminder_GetReminderForChannelIDByIDOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, reminder := range testReminders() {
		mock.ExpectQuery("SELECT (.*) FROM `reminders`").
			WithArgs(
				"!abcdefghijklmop",
				reminder.ID,
			).
			WillReturnError(errors.New("test error"))

		c := &Channel{}
		c.ID = reminder.ChannelID
		_, err := db.GetReminderForChannelIDByID("!abcdefghijklmop", int(reminder.ID))
		require.Error(t, err)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestReminder_GetPendingRemindersOnSuccess(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, reminder := range testReminders() {
		mock.ExpectQuery("SELECT (.*) FROM `reminders`").
			WithArgs(
				reminder.Channel.ChannelIdentifier,
				true,
			).
			WillReturnRows(rowsForReminders([]*Reminder{reminder}))

		c := &Channel{}
		c.ID = reminder.ChannelID
		newReminders, err := db.GetPendingReminders(c)
		require.NoError(t, err)

		require.Equal(t, 1, len(newReminders))
		assert.Equal(reminder.RemindTime, newReminders[0].RemindTime)
		assert.Equal(reminder.Message, newReminders[0].Message)
		assert.Equal(reminder.Active, newReminders[0].Active)
		assert.Equal(reminder.RepeatInterval, newReminders[0].RepeatInterval)
		assert.Equal(reminder.RepeatMax, newReminders[0].RepeatMax)
		assert.Equal(reminder.Repeated, newReminders[0].Repeated)
		assert.Equal(reminder.ChannelID, newReminders[0].ChannelID)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestReminder_GetPendingRemindersOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, reminder := range testReminders() {
		mock.ExpectQuery("SELECT (.*) FROM `reminders`").
			WithArgs(
				reminder.Channel.ChannelIdentifier,
				true,
			).
			WillReturnError(errors.New("test error"))

		c := &Channel{}
		c.ID = reminder.ChannelID
		_, err := db.GetPendingReminders(c)
		assert.Error(err)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestReminder_GetPendingReminderOnSuccess(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	mock.ExpectQuery("SELECT (.*) FROM `reminders`").
		WithArgs(
			int64(1),
			sqlmock.AnyArg(),
		).
		WillReturnRows(rowsForReminders(testReminders()))
	mock.ExpectQuery("SELECT (.*) FROM `channels`").
		WithArgs(
			1,
		).
		WillReturnRows(rowsForChannels([]*Channel{}))

	newReminders, err := db.GetPendingReminder()
	require.NoError(t, err)
	require.Equal(t, len(testReminders()), len(newReminders))

	for _, reminder := range testReminders() {
		found := false
		for _, newReminder := range newReminders {
			if reminder.ID == newReminder.ID {
				found = true
				// RemindTime seems to get fucked up in conversions
				assert.Equal(reminder.Message, newReminder.Message)
				assert.Equal(reminder.Active, newReminder.Active)
				assert.Equal(reminder.RepeatInterval, newReminder.RepeatInterval)
				assert.Equal(reminder.RepeatMax, newReminder.RepeatMax)
				assert.Equal(reminder.Repeated, newReminder.Repeated)
				assert.Equal(reminder.ChannelID, newReminder.ChannelID)
			}
		}
		assert.True(found, "Missing reminder ", reminder.ID)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestReminder_GetPendingReminderOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	mock.ExpectQuery("SELECT (.*) FROM `reminders`").
		WithArgs(
			int64(1),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("test error"))

	_, err := db.GetPendingReminder()
	assert.Error(err)

	assert.NoError(mock.ExpectationsWereMet())
}

func TestReminder_GetMessageFromReminderOnSuccess(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, message := range testMessages() {
		mock.ExpectQuery("SELECT (.*) FROM `messages`").
			WithArgs(
				message.ChannelID,
				message.Type,
			).
			WillReturnRows(rowsForMessages([]*Message{message}))
		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(
				message.ChannelID,
			).
			WillReturnRows(rowsForChannels([]*Channel{}))

		newMsg, err := db.GetMessageFromReminder(*message.ReminderID, message.Type)
		require.NoError(t, err)
		assert.Equal(message.ExternalIdentifier, newMsg.ExternalIdentifier)
		assert.Equal(message.Body, newMsg.Body)
		assert.Equal(message.BodyHTML, newMsg.BodyHTML)
		assert.Equal(message.ResponseToMessage, newMsg.ResponseToMessage)
		assert.Equal(message.Type, newMsg.Type)
		assert.Equal(message.Timestamp, newMsg.Timestamp)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestReminder_GetMessageFromReminderOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, message := range testMessages() {
		mock.ExpectQuery("SELECT (.*) FROM `messages`").
			WithArgs(
				message.ChannelID,
				message.Type,
			).
			WillReturnError(errors.New("test error"))

		_, err := db.GetMessageFromReminder(*message.ReminderID, message.Type)
		assert.Error(err)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func rowsForReminders(reminders []*Reminder) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "remind_time", "message", "active", "repeat_interval", "repeat_max", "repeated", "channel_id"})

	for _, reminder := range reminders {
		rows.AddRow(
			reminder.ID,
			reminder.CreatedAt,
			reminder.UpdatedAt,
			reminder.DeletedAt,
			reminder.RemindTime,
			reminder.Message,
			reminder.Active,
			reminder.RepeatInterval,
			reminder.RepeatMax,
			reminder.Repeated,
			reminder.ChannelID,
		)
	}

	return rows
}

func TestReminder_GetDailyReminderOnSuccess(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, reminder := range testReminders() {
		mock.ExpectQuery("SELECT (.*) FROM `reminders`").
			WithArgs(
				reminder.Channel.ChannelIdentifier,
				sqlmock.AnyArg(),
				true,
			).
			WillReturnRows(rowsForReminders([]*Reminder{reminder}))

		c := &Channel{}
		c.ID = reminder.ChannelID
		newReminders, err := db.GetDailyReminder(c)
		require.NoError(t, err)

		reminderList := *newReminders
		require.Equal(t, 1, len(reminderList))
		assert.Equal(reminder.RemindTime, reminderList[0].RemindTime)
		assert.Equal(reminder.Message, reminderList[0].Message)
		assert.Equal(reminder.Active, reminderList[0].Active)
		assert.Equal(reminder.RepeatInterval, reminderList[0].RepeatInterval)
		assert.Equal(reminder.RepeatMax, reminderList[0].RepeatMax)
		assert.Equal(reminder.Repeated, reminderList[0].Repeated)
		assert.Equal(reminder.ChannelID, reminderList[0].ChannelID)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestReminder_GetDailyReminderOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, reminder := range testReminders() {
		mock.ExpectQuery("SELECT (.*) FROM `reminders`").
			WithArgs(
				reminder.Channel.ChannelIdentifier,
				sqlmock.AnyArg(),
				true,
			).
			WillReturnError(errors.New("test error"))

		c := &Channel{}
		c.ID = reminder.ChannelID
		_, err := db.GetDailyReminder(c)
		assert.Error(err)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func testReminders() []*Reminder {
	reminders := make([]*Reminder, 0)
	reminders = append(reminders, testReminder1())

	return reminders
}

func testReminder1() *Reminder {
	repeated := uint64(1)
	reminder := &Reminder{
		RemindTime:     time.Now(),
		Message:        "abcdegfdgh",
		Active:         true,
		RepeatInterval: 60,
		RepeatMax:      2,
		Repeated:       &repeated,
		ChannelID:      uint(1),
	}
	reminder.ID = 1
	return reminder
}
