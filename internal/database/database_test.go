package database

import (
	"os"
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func testDatabase() (*Database, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	database := Database{
		db: *gormDB,
	}

	return &database, mock
}

func testChannels() []*Channel {
	channels := make([]*Channel, 0)

	channels = append(channels, testChannel1())
	channels = append(channels, testChannel2())

	return channels
}

func testChannel1() *Channel {
	role := roles.RoleAdmin
	c := &Channel{
		Created:           time.Now(),
		ChannelIdentifier: "!abcdefghijklmop",
		UserIdentifier:    "@remindme:matrix.org",
		TimeZone:          "",
		DailyReminder:     nil,
		CalendarSecret:    "abcdefghijklmnop",
		Role:              &role,
	}

	c.ID = 1
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	return c
}

func testChannel2() *Channel {
	role := roles.RoleUser
	c := &Channel{
		Created:           time.Now(),
		ChannelIdentifier: "!123defghijklmop",
		UserIdentifier:    "@remindme2:matrix.org",
		TimeZone:          "Europe/Berlin",
		DailyReminder:     nil,
		CalendarSecret:    "abcdefghijklmnop",
		Role:              &role,
	}

	c.ID = 2
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	return c
}
