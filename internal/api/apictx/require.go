package apictx

import (
	"github.com/gin-gonic/gin"
)

type idInURI struct {
	ID uint `uri:"id" binding:"required"`
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
