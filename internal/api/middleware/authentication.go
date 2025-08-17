package middleware

import (
	"errors"
	"net/http"

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
			for _, value := range values {
				if value == apikey {
					authenticated = true
					break
				}
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
