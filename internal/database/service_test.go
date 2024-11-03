package database_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var logger *slog.Logger
var service database.Service
var gormDB *gorm.DB

func getConnection() string {
	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	return "root:mypass@tcp(" + host + ":3306)/remindme"
}

func getService(config *database.Config) database.Service {
	service, err := database.NewService(
		config,
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

func getLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func TestMain(m *testing.M) {
	logger = getLogger()
	service = getService(&database.Config{
		Connection:    getConnection(),
		LogStatements: true,
	})
	gormDB = getGormDB()

	m.Run()
}

func TestNewService(t *testing.T) {
	_, err := database.NewService(
		&database.Config{
			Connection:    getConnection(),
			LogStatements: true,
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
