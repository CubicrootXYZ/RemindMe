package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WithData responds with the given data and a 200 status code.
func WithData(ctx *gin.Context, data any) {
	resp := DataResponse{
		Status: "success",
		Data:   data,
	}

	ctx.JSON(http.StatusOK, resp)
}
