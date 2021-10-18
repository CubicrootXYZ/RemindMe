package database

import (
	"errors"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannel_AddChannelOnSuccess(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db, mock := testDatabase()

	for _, channel := range testChannels() {

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `channels`").WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			channel.ChannelIdentifier,
			channel.UserIdentifier,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			roles.RoleUser).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(channel.UserIdentifier, channel.ChannelIdentifier).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"}).AddRow(
					channel.ID,
					channel.CreatedAt,
					channel.UpdatedAt,
					channel.DeletedAt,
					channel.Created,
					channel.ChannelIdentifier,
					channel.UserIdentifier,
					channel.TimeZone,
					channel.DailyReminder,
					channel.CalendarSecret,
					channel.Role,
				))

		channelCreated, err := db.AddChannel(channel.UserIdentifier, channel.ChannelIdentifier, roles.RoleUser)
		require.NoError(err)

		assert.Equal(channel.ChannelIdentifier, channelCreated.ChannelIdentifier)
		assert.Equal(channel.UserIdentifier, channelCreated.UserIdentifier)
		assert.Equal(channel.ID, channelCreated.ID)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_AddChannelOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, channel := range testChannels() {

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `channels`").WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			channel.ChannelIdentifier,
			channel.UserIdentifier,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			roles.RoleUser).WillReturnError(errors.New("test error"))
		mock.ExpectRollback()

		_, err := db.AddChannel(channel.UserIdentifier, channel.ChannelIdentifier, roles.RoleUser)
		assert.Error(err)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_TimezoneOnSuccess(t *testing.T) {
	assert := assert.New(t)
	channel := testChannel1()

	zones := []string{"Europe/Berlin", "Africa/Bamako", "America/Mexico_City", "Asia/Jakarta", "Europe/Budapest", "Pacific/Fiji", "US/Samoa"}

	for _, zone := range zones {
		channel.TimeZone = zone

		assert.Equalf(zone, channel.Timezone().String(), "Timezone %s returned as %s", zone, channel.Timezone().String())
	}
}

func TestChannel_TimezoneOnFailure(t *testing.T) {
	assert := assert.New(t)
	channel := testChannel1()

	zones := []string{"Euro/Berlin", "Africa Bamako", "Mexico City", "Asia", "now", "", "1234"}

	for _, zone := range zones {
		channel.TimeZone = zone

		assert.Equalf("UTC", channel.Timezone().String(), "Timezone %s returned as %s", zone, channel.Timezone().String())
	}
}

func TestChannel_GetChannelOnSuccess(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db, mock := testDatabase()

	for _, channel := range testChannels() {

		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(channel.ID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"}).AddRow(
					channel.ID,
					channel.CreatedAt,
					channel.UpdatedAt,
					channel.DeletedAt,
					channel.Created,
					channel.ChannelIdentifier,
					channel.UserIdentifier,
					channel.TimeZone,
					channel.DailyReminder,
					channel.CalendarSecret,
					channel.Role,
				))

		c, err := db.GetChannel(channel.ID)
		require.NoError(err)

		assert.Equal(channel.ChannelIdentifier, c.ChannelIdentifier)
		assert.Equal(channel.UserIdentifier, c.UserIdentifier)
		assert.Equal(channel.ID, c.ID)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_GetChannelOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, channel := range testChannels() {

		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(channel.ID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"}))

		_, err := db.GetChannel(channel.ID)
		assert.Error(err)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_GetChannelByUserIdentifierOnSuccess(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db, mock := testDatabase()

	for _, channel := range testChannels() {

		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(channel.UserIdentifier).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"}).AddRow(
					channel.ID,
					channel.CreatedAt,
					channel.UpdatedAt,
					channel.DeletedAt,
					channel.Created,
					channel.ChannelIdentifier,
					channel.UserIdentifier,
					channel.TimeZone,
					channel.DailyReminder,
					channel.CalendarSecret,
					channel.Role,
				))

		c, err := db.GetChannelByUserIdentifier(channel.UserIdentifier)
		require.NoError(err)

		assert.Equal(channel.ChannelIdentifier, c.ChannelIdentifier)
		assert.Equal(channel.UserIdentifier, c.UserIdentifier)
		assert.Equal(channel.ID, c.ID)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_GetChannelByUserIdentifierOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, channel := range testChannels() {

		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(channel.UserIdentifier).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"}))

		_, err := db.GetChannelByUserIdentifier(channel.UserIdentifier)
		assert.Error(err)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_GetChannelByChannelIdentifierOnSuccess(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db, mock := testDatabase()

	for _, channel := range testChannels() {
		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(channel.ChannelIdentifier).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"}).AddRow(
					channel.ID,
					channel.CreatedAt,
					channel.UpdatedAt,
					channel.DeletedAt,
					channel.Created,
					channel.ChannelIdentifier,
					channel.UserIdentifier,
					channel.TimeZone,
					channel.DailyReminder,
					channel.CalendarSecret,
					channel.Role,
				).AddRow(
					channel.ID+10,
					channel.CreatedAt,
					channel.UpdatedAt,
					channel.DeletedAt,
					channel.Created,
					channel.ChannelIdentifier,
					channel.UserIdentifier,
					channel.TimeZone,
					channel.DailyReminder,
					channel.CalendarSecret,
					channel.Role,
				))

		cs, err := db.GetChannelsByChannelIdentifier(channel.ChannelIdentifier)
		require.NoError(err)

		found := false
		for _, c := range cs {
			if c.ID == channel.ID {
				found = true
				assert.Equal(channel.ChannelIdentifier, c.ChannelIdentifier)
				assert.Equal(channel.UserIdentifier, c.UserIdentifier)
				assert.Equal(channel.ID, c.ID)
			}
		}

		assert.True(found, "Channel not found in returned list: ", channel.ID)

	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_GetChannelByChannelIdentifierOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, channel := range testChannels() {
		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(channel.ChannelIdentifier).
			WillReturnError(errors.New("test error"))

		_, err := db.GetChannelsByChannelIdentifier(channel.ChannelIdentifier)
		assert.Error(err)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_GetChannelByUserAndChannelIdentifierOnSuccess(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db, mock := testDatabase()

	for _, channel := range testChannels() {

		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(channel.UserIdentifier, channel.ChannelIdentifier).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"}).AddRow(
					channel.ID,
					channel.CreatedAt,
					channel.UpdatedAt,
					channel.DeletedAt,
					channel.Created,
					channel.ChannelIdentifier,
					channel.UserIdentifier,
					channel.TimeZone,
					channel.DailyReminder,
					channel.CalendarSecret,
					channel.Role,
				))

		c, err := db.GetChannelByUserAndChannelIdentifier(channel.UserIdentifier, channel.ChannelIdentifier)
		require.NoError(err)

		assert.Equal(channel.ChannelIdentifier, c.ChannelIdentifier)
		assert.Equal(channel.UserIdentifier, c.UserIdentifier)
		assert.Equal(channel.ID, c.ID)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_GetChannelByUserAndChannelIdentifieOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, channel := range testChannels() {

		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(channel.UserIdentifier, channel.ChannelIdentifier).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"}))

		_, err := db.GetChannelByUserAndChannelIdentifier(channel.UserIdentifier, channel.ChannelIdentifier)
		assert.Error(err)
	}

	for _, channel := range testChannels() {

		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(channel.UserIdentifier, channel.ChannelIdentifier).
			WillReturnError(errors.New("test error"))

		_, err := db.GetChannelByUserAndChannelIdentifier(channel.UserIdentifier, channel.ChannelIdentifier)
		assert.Error(err)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_GetChannelListOnSuccess(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db, mock := testDatabase()

	response := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"})

	for _, channel := range testChannels() {
		response.AddRow(
			channel.ID,
			channel.CreatedAt,
			channel.UpdatedAt,
			channel.DeletedAt,
			channel.Created,
			channel.ChannelIdentifier,
			channel.UserIdentifier,
			channel.TimeZone,
			channel.DailyReminder,
			channel.CalendarSecret,
			channel.Role,
		)
	}

	mock.ExpectQuery("SELECT (.*) FROM `channels`").WillReturnRows(response)

	cs, err := db.GetChannelList()
	require.NoError(err)

	for _, channel := range testChannels() {
		found := false
		for _, c := range cs {
			if channel.ID == c.ID {
				found = true
				assert.Equal(channel.ChannelIdentifier, c.ChannelIdentifier)
				assert.Equal(channel.UserIdentifier, c.UserIdentifier)
				assert.Equal(channel.ID, c.ID)
			}

		}
		assert.True(found, "Channel ID not found in response: ", channel.ID)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_GetChannelListOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	mock.ExpectQuery("SELECT (.*) FROM `channels`").WillReturnError(errors.New("test error"))

	_, err := db.GetChannelList()
	assert.Error(err)

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_ChannelCountOnSuccess(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db, mock := testDatabase()

	mock.ExpectQuery("SELECT (.*) FROM `channels`").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(2))

	count, err := db.ChannelCount()
	require.NoError(err)
	assert.Equalf(int64(2), count, "Received %d channels but it should be %d", count, 2)

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_ChannelCountOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	mock.ExpectQuery("SELECT (.*) FROM `channels`").WillReturnError(errors.New("test error"))

	_, err := db.ChannelCount()
	assert.Error(err)

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_UpdateChannelOnSuccess(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db, mock := testDatabase()

	remindTime := uint(6)
	role := roles.RoleAdmin

	for _, channel := range testChannels() {
		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(
				channel.ID,
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"}).AddRow(
					channel.ID,
					channel.CreatedAt,
					channel.UpdatedAt,
					channel.DeletedAt,
					channel.Created,
					channel.ChannelIdentifier,
					channel.UserIdentifier,
					channel.TimeZone,
					channel.DailyReminder,
					channel.CalendarSecret,
					channel.Role,
				))

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `channels`").WithArgs(
			channel.CreatedAt,
			sqlmock.AnyArg(),
			channel.DeletedAt,
			channel.Created,
			channel.ChannelIdentifier,
			channel.UserIdentifier,
			"Europe/Berlin",
			&remindTime,
			channel.CalendarSecret,
			&role,
			channel.ID,
		).WillReturnResult(sqlmock.NewResult(int64(channel.ID), 1))
		mock.ExpectCommit()

		c, err := db.UpdateChannel(channel.ID, "Europe/Berlin", &remindTime, &role)
		require.NoError(err)

		assert.Equal(channel.ID, c.ID)
		assert.Equal("Europe/Berlin", c.TimeZone)
		assert.Equal(&remindTime, c.DailyReminder)
		assert.Equal(&role, c.Role)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestChannel_UpdateChannelOnFailure(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db, mock := testDatabase()

	remindTime := uint(6)
	role := roles.RoleAdmin

	for _, channel := range testChannels() {
		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(
				channel.ID,
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"}).AddRow(
					channel.ID,
					channel.CreatedAt,
					channel.UpdatedAt,
					channel.DeletedAt,
					channel.Created,
					channel.ChannelIdentifier,
					channel.UserIdentifier,
					channel.TimeZone,
					channel.DailyReminder,
					channel.CalendarSecret,
					channel.Role,
				))

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `channels`").WithArgs(
			channel.CreatedAt,
			sqlmock.AnyArg(),
			channel.DeletedAt,
			channel.Created,
			channel.ChannelIdentifier,
			channel.UserIdentifier,
			"Europe/Berlin",
			&remindTime,
			channel.CalendarSecret,
			&role,
			channel.ID,
		).WillReturnError(errors.New("test error"))

		_, err := db.UpdateChannel(channel.ID, "Europe/Berlin", &remindTime, &role)
		require.Error(err)
	}

	for _, channel := range testChannels() {
		mock.ExpectQuery("SELECT (.*) FROM `channels`").
			WithArgs(
				channel.ID,
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "created", "channel_identifier", "user_identifier", "time_zone", "daily_reminder", "calendar_secret", "role"}))

		_, err := db.UpdateChannel(channel.ID, "Europe/Berlin", &remindTime, &role)
		require.Error(err)
	}

	assert.NoError(mock.ExpectationsWereMet())
}
