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
	bl.ID = 1
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
