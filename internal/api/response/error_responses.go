package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AbortWithInternalServerError aborts request.
func AbortWithInternalServerError(ctx *gin.Context) {
	response := MessageErrorResponse{
		Status:  "error",
		Message: "Internal Server Error",
	}
	ctx.JSON(http.StatusInternalServerError, response)
	ctx.Abort()
}

// AbortWithNotFoundError aborts request.
func AbortWithNotFoundError(ctx *gin.Context) {
	response := MessageErrorResponse{
		Status:  "error",
		Message: "Not Found",
	}
	ctx.JSON(http.StatusNotFound, response)
	ctx.Abort()
}
