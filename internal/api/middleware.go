package api

import (
	"net/http"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"github.com/gin-gonic/gin"
)

type idInURI struct {
	ID uint `uri:"id" binding:"required"`
}

// RequireAPIkey is a middleware that requires the given api key to be present in the headers authorization field or as a url parameter.
func RequireAPIkey(apikey string) gin.HandlerFunc {
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
			response := types.MessageErrorResponse{
				Message: "Unauthenticated",
				Status:  "error",
			}
			ctx.JSON(http.StatusUnauthorized, response)
			ctx.AbortWithError(http.StatusForbidden, errors.ErrMissingApiKey)
			return
		}
	}
}

// RequireCalendarSecret is a middleware that requires the calendar secret to be set as a query parameter.
func RequireCalendarSecret() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		log.Info(token)

		ctx.Set("token", token)
	}
}

// RequireIDInURI returns a Gin middleware which requires an ID to be supplied in the URI of the request.
func RequireIDInURI() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestModel idInURI

		if err := ctx.BindUri(&requestModel); err != nil {
			return
		}

		ctx.Set("id", requestModel.ID)
	}
}
