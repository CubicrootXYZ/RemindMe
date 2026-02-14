package middleware

import (
	"errors"
	"net/http"
	"slices"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/response"
	"github.com/gin-gonic/gin"
)

// APIKeyAuth is a middleware that enforces the API key to be set
// as "Authorization" header.
func APIKeyAuth(apikey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headers := ctx.Request.Header

		authenticated := false

		if values, ok := headers["Authorization"]; ok {
			if slices.Contains(values, apikey) {
				authenticated = true
			}
		}

		if !authenticated {
			response := response.MessageErrorResponse{
				Message: "Unauthenticated",
				Status:  "error",
			}
			ctx.JSON(http.StatusUnauthorized, response)
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("missing api key"))

			return
		}
	}
}
