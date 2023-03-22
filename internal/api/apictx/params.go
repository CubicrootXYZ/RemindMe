package apictx

import "github.com/gin-gonic/gin"

// GetUintFromContext returns the ID from the context.
func GetUintFromContext(ctx *gin.Context, name string) (uint, bool) {
	valueRaw, ok := ctx.Get(name)
	if !ok {
		return 0, false
	}

	value, ok := valueRaw.(uint)
	if !ok {
		return 0, false
	}

	return value, true
}
