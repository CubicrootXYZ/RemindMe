package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvent_IsEventKnownOnSuccess(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db, mock := testDatabase()

	for _, event := range testEvents() {
		mock.ExpectQuery("SELECT (.*) FROM `events`").WillReturnRows(rowsForEvents([]*Event{event}))

		exists, err := db.IsEventKnown(event.ExternalIdentifier)

		require.NoError(err)
		assert.Truef(exists, "Event %d does not exist", event.ID)
	}

	for _, event := range testEvents() {
		mock.ExpectQuery("SELECT (.*) FROM `events`").WillReturnRows(rowsForEvents([]*Event{}))

		exists, err := db.IsEventKnown(event.ExternalIdentifier)

		require.NoError(err)
		assert.Falsef(exists, "Event %d should not exist", event.ID)
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestEvent_IsEventKnownOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, event := range testEvents() {
		mock.ExpectQuery("SELECT (.*) FROM `events`").WillReturnError(errors.New("test error"))

		_, err := db.IsEventKnown(event.ExternalIdentifier)

		assert.Error(err)
	}
	assert.NoError(mock.ExpectationsWereMet())
}

func TestEvent_AddEventOnSuccess(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, event := range testEvents() {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `events`").WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			event.ChannelID,
			event.Timestamp,
			event.ExternalIdentifier,
			event.EventType,
			event.EventSubType,
			event.AdditionalInfo).
			WillReturnResult(sqlmock.NewResult(int64(event.ID), 1))
		mock.ExpectCommit()

		newEvent, err := db.AddEvent(event)
		require.NoError(t, err)

		assert.Equal(event.ChannelID, newEvent.ChannelID)
		assert.Equal(event.Timestamp, newEvent.Timestamp)
		assert.Equal(event.ExternalIdentifier, newEvent.ExternalIdentifier)
		assert.Equal(event.EventType, newEvent.EventType)
		assert.Equal(event.EventSubType, newEvent.EventSubType)
		assert.Equal(event.AdditionalInfo, newEvent.AdditionalInfo)
	}
	assert.NoError(mock.ExpectationsWereMet())
}

func TestEvent_AddEventOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, event := range testEvents() {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `events`").WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			event.ChannelID,
			event.Timestamp,
			event.ExternalIdentifier,
			event.EventType,
			event.EventSubType,
			event.AdditionalInfo).
			WillReturnError(errors.New("test error"))
		mock.ExpectRollback()

		_, err := db.AddEvent(event)
		assert.Error(err)
	}
	assert.NoError(mock.ExpectationsWereMet())
}

func testEvents() []*Event {
	events := make([]*Event, 0)
	events = append(events, testEvent1())

	return events
}

func testEvent1() *Event {
	id := uint(1)
	return &Event{
		ChannelID:          &id,
		Timestamp:          12345,
		ExternalIdentifier: "abcdefg",
		EventType:          EventTypeMembership,
		EventSubType:       "test",
		AdditionalInfo:     "",
	}
}

func rowsForEvents(events []*Event) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "channel_id", "timestamp", "external_identifier", "event_type", "event_sub_type", "additional_info"})

	for _, event := range events {
		rows.AddRow(
			event.ID,
			event.CreatedAt,
			event.UpdatedAt,
			event.DeletedAt,
			event.ChannelID,
			event.Timestamp,
			event.ExternalIdentifier,
			event.EventType,
			event.EventSubType,
			event.AdditionalInfo,
		)
	}

	return rows
}
