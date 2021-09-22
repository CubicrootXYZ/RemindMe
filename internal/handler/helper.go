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

// getUintFromContext returns the ID from the context
func getUintFromContext(ctx *gin.Context, name string) (uint, error) {
	id, ok := ctx.MustGet(name).(uint)
	if !ok {
		return 0, errors.ErrMissingID
	}

	return id, nil
}

// getStringFromContext returns the ID from the context
func getStringFromContext(ctx *gin.Context, name string) (string, error) {
	text, ok := ctx.MustGet(name).(string)
	if !ok {
		return "", errors.ErrMissingID // TODO make separate errors
	}

	return text, nil
}
