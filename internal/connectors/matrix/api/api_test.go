package api_test

import (
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/api"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
)

func testServer(t *testing.T) (*database.MockService, *matrixdb.MockService, *httptest.Server) { //nolint:unparam
	t.Helper()
	db := database.NewMockService(t)
	matrixDB := matrixdb.NewMockService(t)

	api := api.New(&api.Config{
		Database:            db,
		MatrixDB:            matrixDB,
		DefaultAuthProvider: func(_ *gin.Context) {},
	}, slog.Default())

	r := gin.New()

	err := api.RegisterRoutes(r)
	if err != nil {
		panic(err)
	}

	return db, matrixDB, httptest.NewServer(r)
}
