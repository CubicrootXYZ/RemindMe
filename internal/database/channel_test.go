package database

import (
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestChannel_AddChannel(t *testing.T) {
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
