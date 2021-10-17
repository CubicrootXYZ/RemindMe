package database

import (
	"errors"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestChannel_AddChannelOnSuccess(t *testing.T) {
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

		_, err := db.AddChannel(channel.UserIdentifier, channel.ChannelIdentifier, roles.RoleUser)
		assert.NoError(err)
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
