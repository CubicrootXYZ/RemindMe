package tests

import (
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

	db, err := database.Create(configuration.Database{
		Connection: "root:mypass@tcp(database:3306)/remindme",
	}, true)
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(mysql.Open("root:mypass@tcp(database:3306)/remindme?parseTime=True"), &gorm.Config{})
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
	err = gormDatabase.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	if err != nil {
		panic(err)
	}

	_, err = db.AddChannel(
		"testuser@example.com",
		"!123456789",
		roles.RoleUser,
	)
	if err != nil {
		panic(err)
	}

	_, err = db.AddChannel(
		"testuser2@example.com",
		"!abcdefghij",
		roles.RoleAdmin,
	)
	if err != nil {
		panic(err)
	}
}
