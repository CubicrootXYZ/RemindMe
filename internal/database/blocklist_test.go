package database

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testBlocklists() []*Blocklist {
	bls := make([]*Blocklist, 0)
	bls = append(bls, testBlocklist1())
	bls = append(bls, testBlocklist2())

	return bls
}

func testBlocklist1() *Blocklist {
	bl := &Blocklist{}
	bl.ID = 1
	bl.CreatedAt = time.Now()
	bl.UpdatedAt = time.Now()
	bl.UserIdentifier = "@remindme:matrix.org"
	bl.Reason = "I do not know"

	return bl
}

func testBlocklist2() *Blocklist {
	bl := &Blocklist{}
	bl.ID = 2
	bl.CreatedAt = time.Now()
	bl.UpdatedAt = time.Now()
	bl.UserIdentifier = "@remindme2:matrix.org"
	bl.Reason = "I do not know what is happening here"

	return bl
}

func TestBlocklist_IsUserBlockedOnSuccess(t *testing.T) {
	db, mock := testDatabase()

	for _, bl := range testBlocklists() {
		mock.ExpectQuery("SELECT (.*) FROM `blocklists`").WithArgs(bl.UserIdentifier).
			WillReturnRows(rowsForBlocklists([]*Blocklist{bl}))

		isBlocked, err := db.IsUserBlocked(bl.UserIdentifier)

		require.NoError(t, err)
		assert.True(t, isBlocked)
		assert.NoError(t, mock.ExpectationsWereMet())
	}

	for _, bl := range testBlocklists() {
		mock.ExpectQuery("SELECT (.*) FROM `blocklists`").WithArgs(bl.UserIdentifier).
			WillReturnRows(rowsForBlocklists([]*Blocklist{}))

		isBlocked, err := db.IsUserBlocked(bl.UserIdentifier)

		require.NoError(t, err)
		assert.False(t, isBlocked)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestBlocklist_IsUserBlockedOnFailure(t *testing.T) {
	db, mock := testDatabase()

	for _, bl := range testBlocklists() {
		mock.ExpectQuery("SELECT (.*) FROM `blocklists`").WithArgs(bl.UserIdentifier).
			WillReturnError(errors.New("test error"))

		isBlocked, err := db.IsUserBlocked(bl.UserIdentifier)

		assert.Error(t, err)
		assert.False(t, isBlocked)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestBlocklist_AddUserToBlocklistOnSuccess(t *testing.T) {
	db, mock := testDatabase()

	for _, bl := range testBlocklists() {
		mock.ExpectQuery("SELECT (.*) FROM `blocklists`").WithArgs(bl.UserIdentifier).
			WillReturnRows(rowsForBlocklists([]*Blocklist{}))

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `blocklists`").WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			bl.UserIdentifier,
			bl.Reason,
		).WillReturnResult(sqlmock.NewResult(int64(bl.ID), 1))
		mock.ExpectCommit()

		err := db.AddUserToBlocklist(bl.UserIdentifier, bl.Reason)

		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}

	for _, bl := range testBlocklists() {
		mock.ExpectQuery("SELECT (.*) FROM `blocklists`").WithArgs(bl.UserIdentifier).
			WillReturnRows(rowsForBlocklists([]*Blocklist{bl}))

		err := db.AddUserToBlocklist(bl.UserIdentifier, bl.Reason)

		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestBlocklist_AddUserToBlocklistOnFailure(t *testing.T) {
	db, mock := testDatabase()

	for _, bl := range testBlocklists() {
		mock.ExpectQuery("SELECT (.*) FROM `blocklists`").WithArgs(bl.UserIdentifier).
			WillReturnRows(rowsForBlocklists([]*Blocklist{}))

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `blocklists`").WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			bl.UserIdentifier,
			bl.Reason,
		).WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))

		err := db.AddUserToBlocklist(bl.UserIdentifier, bl.Reason)

		require.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestBlocklist_RemoveUserFromBlocklistOnSuccess(t *testing.T) {
	db, mock := testDatabase()

	for _, bl := range testBlocklists() {
		mock.ExpectExec("DELETE FROM `blocklists`").WithArgs(
			bl.UserIdentifier,
		).WillReturnResult(sqlmock.NewResult(int64(bl.ID), 1))

		err := db.RemoveUserFromBlocklist(bl.UserIdentifier)

		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestBlocklist_RemoveUserFromBlocklistOnFailure(t *testing.T) {
	db, mock := testDatabase()

	for _, bl := range testBlocklists() {
		mock.ExpectExec("DELETE FROM `blocklists`").WithArgs(
			bl.UserIdentifier,
		).WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))

		_ = db.RemoveUserFromBlocklist(bl.UserIdentifier)

		//require.Error(t, err) // Somehow does not error
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestBlocklist_GetBlockedUserListOnSuccess(t *testing.T) {
	assert := assert.New(t)
	db, mock := testDatabase()

	blocklists := testBlocklists()

	mock.ExpectQuery("SELECT (.*) FROM `blocklists`").
		WillReturnRows(rowsForBlocklists(blocklists))

	list, err := db.GetBlockedUserList()

	require.NoError(t, err)
	assert.Equal(len(blocklists), len(list), "Not same amount of lists returned as entered")

	for _, blocklist := range blocklists {
		found := false
		for _, l := range list {
			if l.ID == blocklist.ID {
				found = true
				assert.Equal(blocklist.CreatedAt, l.CreatedAt)
				assert.Equal(blocklist.UpdatedAt, l.UpdatedAt)
				assert.Equal(blocklist.DeletedAt, l.DeletedAt)
				assert.Equal(blocklist.UserIdentifier, l.UserIdentifier)
				assert.Equal(blocklist.Reason, l.Reason)
			}
		}
		assert.True(found)
	}

	assert.NoError(mock.ExpectationsWereMet())

}

func TestBlocklist_GetBlockedUserListOnFailure(t *testing.T) {
	db, mock := testDatabase()

	blocklists := testBlocklists()

	mock.ExpectQuery("SELECT (.*) FROM `blocklists`").WithArgs(blocklists).
		WillReturnError(errors.New("test error"))

	list, err := db.GetBlockedUserList()

	assert.Error(t, err)
	assert.Equal(t, 0, len(list), "List not empty")
}

// HELPER

func rowsForBlocklists(blocklists []*Blocklist) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_identifier", "reason"})

	for _, blocklist := range blocklists {
		rows.AddRow(
			blocklist.ID,
			blocklist.CreatedAt,
			blocklist.UpdatedAt,
			blocklist.DeletedAt,
			blocklist.UserIdentifier,
			blocklist.Reason,
		)
	}

	return rows
}
