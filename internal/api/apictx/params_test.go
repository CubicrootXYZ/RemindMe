package apictx_test

import (
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/apictx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestGetUintFromContext(t *testing.T) {
	ctx := &gin.Context{}
	ctx.Set("test", uint(123))

	value, ok := apictx.GetUintFromContext(ctx, "test")
	require.True(t, ok)
	assert.Equal(t, uint(123), value)
}

func TestGetUintFromContextWithWrongType(t *testing.T) {
	ctx := &gin.Context{}
	ctx.Set("test", 123)

	_, ok := apictx.GetUintFromContext(ctx, "test")
	require.False(t, ok)
}

func TestGetUintFromContextWithNotSet(t *testing.T) {
	ctx := &gin.Context{}

	_, ok := apictx.GetUintFromContext(ctx, "test")
	require.False(t, ok)
}
