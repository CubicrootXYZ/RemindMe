package api_test

import (
	"net/http/httptest"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/api"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func testServer(ctrl *gomock.Controller) (*database.MockService, *matrixdb.MockService, *httptest.Server) { //nolint:unparam
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)

	api := api.New(&api.Config{
		Database:            db,
		MatrixDB:            matrixDB,
		DefaultAuthProvider: func(_ *gin.Context) {},
	}, gologger.New(gologger.LogLevelDebug, 0))

	r := gin.New()
	err := api.RegisterRoutes(r)
	if err != nil {
		panic(err)
	}

	return db, matrixDB, httptest.NewServer(r)
}
