package handler

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"github.com/gin-gonic/gin"
)

// abort aborts a handler execution
func abort(ctx *gin.Context, statusCode int, message ResponseMessage, err error) {
	response := types.MessageErrorResponse{
		Status:  "error",
		Message: string(message),
	}
	ctx.JSON(statusCode, response)
	ctx.AbortWithError(statusCode, err)
}

// getIDfromContext returns the ID from the context
func getIDfromContext(ctx *gin.Context) (uint, error) {
	id, ok := ctx.MustGet("id").(uint)
	if !ok {
		return 0, errors.ErrMissingID
	}

	return id, nil
}
