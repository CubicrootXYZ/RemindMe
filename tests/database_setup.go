package tests

import (
	"os"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDatabase *database.Database
var gormDatabase *gorm.DB

func NewDatabase() *database.Database {
	if testDatabase != nil {
		return testDatabase
	}

	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	db, err := database.Create(configuration.Database{
		Connection: "root:mypass@tcp(" + host + ":3306)/remindme",
	}, true)
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(mysql.Open("root:mypass@tcp("+host+":3306)/remindme?parseTime=True"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	gormDatabase = gormDB

	populateFixtures(db)
	testDatabase = db

	return testDatabase
}

func populateFixtures(db *database.Database) {
	err := gormDatabase.Exec("SET FOREIGN_KEY_CHECKS = 0").Error
	if err != nil {
		panic(err)
	}
	err = gormDatabase.Exec("TRUNCATE channels").Error
	if err != nil {
		panic(err)
	}
	err = gormDatabase.Exec("TRUNCATE third_party_resources").Error
	if err != nil {
		panic(err)
	}
	err = gormDatabase.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	if err != nil {
		panic(err)
	}

	// Channel 1
	_, err = db.AddChannel(
		"testuser@example.com",
		"!123456789",
		roles.RoleUser,
	)
	if err != nil {
		panic(err)
	}

	// Channel 2
	c, err := db.AddChannel(
		"testuser2@example.com",
		"!abcdefghij",
		roles.RoleAdmin,
	)
	if err != nil {
		panic(err)
	}

	r := roles.RoleAdmin
	_, err = db.UpdateChannel(c.ID, "Berlin", nil, &r)
	if err != nil {
		panic(err)
	}

	// Channel 3 (to delete)
	_, err = db.AddChannel(
		"testuser3@example.com",
		"!123456789abcde",
		roles.RoleUser,
	)
	if err != nil {
		panic(err)
	}
}
