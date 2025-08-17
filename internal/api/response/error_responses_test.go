package response_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAbortWithInternalServerError(t *testing.T) {
	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		response.AbortWithInternalServerError(ctx)
	})
	server := httptest.NewServer(r)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"message":"Internal Server Error","status":"error"}`, string(body))
}

func TestAbortWithNotFoundError(t *testing.T) {
	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		response.AbortWithNotFoundError(ctx)
	})
	server := httptest.NewServer(r)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL+"/",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"message":"Not Found","status":"error"}`, string(body))
}
