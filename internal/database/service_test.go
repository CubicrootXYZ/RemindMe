package database_test

import (
	"os"
	"testing"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var logger gologger.Logger
var service database.Service
var gormDB *gorm.DB

func getConnection() string {
	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	return "root:mypass@tcp(" + host + ":3306)/remindme"
}

func getService() database.Service {
	service, err := database.NewService(
		&database.Config{
			Connection: getConnection(),
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	return service
}

func getGormDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(getConnection()+"?parseTime=True"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func getLogger() gologger.Logger {
	return gologger.New(gologger.LogLevelDebug, 0)
}

func TestMain(m *testing.M) {
	logger = getLogger()
	service = getService()
	gormDB = getGormDB()

	m.Run()
}

func TestNewService(t *testing.T) {
	_, err := database.NewService(
		&database.Config{
			Connection: getConnection(),
		},
		logger,
	)
	require.NoError(t, err)
}

func TestNewServiceWithNoConnection(t *testing.T) {
	_, err := database.NewService(
		&database.Config{},
		logger,
	)
	require.Error(t, err)
}

func TestNewServiceWithNoConfig(t *testing.T) {
	_, err := database.NewService(
		nil,
		logger,
	)
	require.Error(t, err)
}
