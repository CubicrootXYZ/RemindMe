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

func TestWithDataWithString(t *testing.T) {
	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		response.WithData(ctx, "hello world")
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"status":"success","data":"hello world"}`, string(body))
}

func TestWithDataWithMap(t *testing.T) {
	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		response.WithData(ctx, map[string]interface{}{
			"key":      "value",
			"key2":     1,
			"k e y 3 ": []uint64{1, 2, 3},
		})
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"status":"success","data":{"k e y 3 ":[1,2,3],"key":"value","key2":1}}`, string(body))
}

func TestWithDataWithNil(t *testing.T) {
	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		response.WithData(ctx, nil)
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"status":"success","data":null}`, string(body))
}
