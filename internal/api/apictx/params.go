package apictx

import "github.com/gin-gonic/gin"

// GetUintFromContext returns the ID from the context.
func GetUintFromContext(ctx *gin.Context, name string) (uint, bool) {
	// TODO test
	id, ok := ctx.MustGet(name).(uint)
	if !ok {
		return 0, false
	}

	return id, true
}
