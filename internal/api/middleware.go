package api

import (
	"net/http"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"github.com/gin-gonic/gin"
)

type idInURI struct {
	ID uint `uri:"id" binding:"required"`
}

type stringIDInURI struct {
	ID string `uri:"id" binding:"required"`
}

// RequireAPIkey is a middleware that requires the given api key to be present in the headers authorization field.
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
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.ErrMissingAPIKey)
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

// RequireIDInURI returns a Gin middleware which requires an ID of the type uint to be supplied in the URI of the request.
func RequireIDInURI() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestModel idInURI

		if err := ctx.BindUri(&requestModel); err != nil {
			return
		}

		ctx.Set("id", requestModel.ID)
	}
}

// RequireStringIDInURI returns a Gin middleware which requires an ID of the type string to be supplied in the URI of the request.
func RequireStringIDInURI() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestModel stringIDInURI

		if err := ctx.BindUri(&requestModel); err != nil {
			return
		}

		ctx.Set("id", requestModel.ID)
	}
}

// Logger is a generic logger for gin
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		log.UInfo("New request",
			"path", param.Path,
			"status_code", param.StatusCode,
			"method", param.Method,
			"timestamp", param.TimeStamp.Unix(),
			"duration", param.Latency.Seconds(),
			"error", param.ErrorMessage,
		)

	}
}
