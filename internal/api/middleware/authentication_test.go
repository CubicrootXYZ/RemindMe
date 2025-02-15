package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIKeyAuth(t *testing.T) {
	r := gin.New()
	r.Use(middleware.APIKeyAuth("123"))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	svr := httptest.NewServer(r)
	defer svr.Close()

	t.Run("happy case", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, svr.URL+"/ping", nil)
		require.NoError(t, err)

		req.Header.Add("Authorization", "123")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.JSONEq(t, "{\"message\":\"pong\"}", string(body))
	})

	t.Run("wrong key", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, svr.URL+"/ping", nil)
		require.NoError(t, err)

		req.Header.Add("Authorization", "1234")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.JSONEq(t, "{\"message\":\"Unauthenticated\",\"status\":\"error\"}", string(body))
	})

	t.Run("no key", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, svr.URL+"/ping", nil)
		require.NoError(t, err)

		req.Header.Add("Authorization", "1234")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.JSONEq(t, "{\"message\":\"Unauthenticated\",\"status\":\"error\"}", string(body))
	})
}
