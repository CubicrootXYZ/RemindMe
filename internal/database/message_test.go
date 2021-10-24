package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage_GetMessageByExternalIDOnSuccess(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, message := range testMessages() {
		mock.ExpectQuery("SELECT (.*) FROM `messages`").
			WithArgs(
				message.ExternalIdentifier,
			).
			WillReturnRows(rowsForMessages([]*Message{message}))
		mock.ExpectQuery("SELECT (.*) FROM `reminders`").
			WithArgs(
				message.ReminderID,
			).
			WillReturnRows(rowsForReminders([]*Reminder{}))

		newMsg, err := db.GetMessageByExternalID(message.ExternalIdentifier)

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

func TestMessage_GetMessageByExternalIDOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, message := range testMessages() {
		mock.ExpectQuery("SELECT (.*) FROM `messages`").
			WithArgs(
				message.ExternalIdentifier,
			).
			WillReturnRows(rowsForMessages([]*Message{}))

		_, err := db.GetMessageByExternalID(message.ExternalIdentifier)

		assert.Error(err)
	}

	for _, message := range testMessages() {
		mock.ExpectQuery("SELECT (.*) FROM `messages`").
			WithArgs(
				message.ExternalIdentifier,
			).
			WillReturnError(errors.New("test error"))

		_, err := db.GetMessageByExternalID(message.ExternalIdentifier)

		assert.Error(err)
	}
	assert.NoError(mock.ExpectationsWereMet())
}

func TestMessage_GetMessagesByReminderIDOnSuccess(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, message := range testMessages() {
		mock.ExpectQuery("SELECT (.*) FROM `messages`").
			WithArgs(
				message.ReminderID,
			).
			WillReturnRows(rowsForMessages([]*Message{message}))

		newMsgs, err := db.GetMessagesByReminderID(*message.ReminderID)

		require.NoError(t, err)
		require.Equal(t, 1, len(newMsgs))
		assert.Equal(message.ExternalIdentifier, newMsgs[0].ExternalIdentifier)
		assert.Equal(message.Body, newMsgs[0].Body)
		assert.Equal(message.BodyHTML, newMsgs[0].BodyHTML)
		assert.Equal(message.ResponseToMessage, newMsgs[0].ResponseToMessage)
		assert.Equal(message.Type, newMsgs[0].Type)
		assert.Equal(message.Timestamp, newMsgs[0].Timestamp)
	}

	for _, message := range testMessages() {
		mock.ExpectQuery("SELECT (.*) FROM `messages`").
			WithArgs(
				message.ReminderID,
			).
			WillReturnRows(rowsForMessages([]*Message{}))

		newMsgs, err := db.GetMessagesByReminderID(*message.ReminderID)

		require.NoError(t, err)
		require.Equal(t, 0, len(newMsgs))
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestMessage_GetMessagesByReminderIDOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, message := range testMessages() {
		mock.ExpectQuery("SELECT (.*) FROM `messages`").
			WithArgs(
				message.ReminderID,
			).
			WillReturnError(errors.New("test error"))

		_, err := db.GetMessagesByReminderID(*message.ReminderID)

		assert.Error(err)
	}
	assert.NoError(mock.ExpectationsWereMet())
}

func testMessages() []*Message {
	messages := make([]*Message, 0)
	messages = append(messages, testMessage1())

	return messages
}

func TestMessage_GetLastMessageByTypeOnSuccess(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, message := range testMessages() {
		mock.ExpectQuery("SELECT (.*) FROM `messages`").
			WithArgs(
				message.ChannelID,
				message.Type,
			).
			WillReturnRows(rowsForMessages([]*Message{message}))

		c := &Channel{}
		c.ID = message.ChannelID
		newMsg, err := db.GetLastMessageByType(message.Type, c)

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

func TestMessage_GetLastMessageByTypeOnFailure(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	for _, message := range testMessages() {
		mock.ExpectQuery("SELECT (.*) FROM `messages`").
			WithArgs(
				message.ChannelID,
				message.Type,
			).
			WillReturnRows(rowsForMessages([]*Message{}))

		c := &Channel{}
		c.ID = message.ChannelID
		_, err := db.GetLastMessageByType(message.Type, c)

		assert.Error(err)
	}

	for _, message := range testMessages() {
		mock.ExpectQuery("SELECT (.*) FROM `messages`").
			WithArgs(
				message.ChannelID,
				message.Type,
			).
			WillReturnError(errors.New("test error"))

		c := &Channel{}
		c.ID = message.ChannelID
		_, err := db.GetLastMessageByType(message.Type, c)

		assert.Error(err)
	}
	assert.NoError(mock.ExpectationsWereMet())
}

func testMessage1() *Message {
	reminderID := uint(45)
	return &Message{
		Body:               "abcde",
		BodyHTML:           "<br>abcde",
		ResponseToMessage:  "hkdjviuz",
		Type:               MessageTypeDailyReminderUpdate,
		ChannelID:          uint(45),
		Timestamp:          12345,
		ExternalIdentifier: "!gjdlaspfoidsj",
		ReminderID:         &reminderID,
	}
}

func rowsForMessages(messages []*Message) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "body_html", "reminder_id", "response_to_message", "type", "channel_id", "timestamp", "external_identifier"})

	for _, message := range messages {
		rows.AddRow(
			message.ID,
			message.CreatedAt,
			message.UpdatedAt,
			message.DeletedAt,
			message.Body,
			message.BodyHTML,
			message.ReminderID,
			message.ResponseToMessage,
			message.Type,
			message.ChannelID,
			message.Timestamp,
			message.ExternalIdentifier,
		)
	}

	return rows
}
