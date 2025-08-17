package apictx

import (
	"github.com/gin-gonic/gin"
)

type idInURI struct {
	ID uint `binding:"required" uri:"id"`
}

// RequireIDInURI returns a Gin middleware which requires an ID of the type uint to be supplied in the URI of the request.
func RequireIDInURI() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestModel idInURI

		err := ctx.BindUri(&requestModel)
		if err != nil {
			return
		}

		ctx.Set("id", requestModel.ID)
	}
}
