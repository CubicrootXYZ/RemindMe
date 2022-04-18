package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireAPIkey(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// test request, must instantiate a request first
	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header), // if you need to test headers
	}
	req.Header.Add("Authorization", "abcdefg")
	c.Request = req

	test := RequireAPIkey("abcdefg")
	test(c)

	assert.Equal(t, w.Code, 200)
	require.Equal(t, 0, len(c.Errors))
}

func TestRequireAPIkey_SpecialCharacters(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// test request, must instantiate a request first
	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header), // if you need to test headers
	}
	req.Header.Add("Authorization", " _d9_/()43ß2#*é")
	c.Request = req

	test := RequireAPIkey(" _d9_/()43ß2#*é")
	test(c)

	assert.Equal(t, w.Code, 200)
	require.Equal(t, 0, len(c.Errors))
}

func TestRequireAPIkey_Empty(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// test request, must instantiate a request first
	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header), // if you need to test headers
	}
	req.Header.Add("Authorization", "")
	c.Request = req

	test := RequireAPIkey("abcdefg")
	test(c)

	assert.Equal(t, w.Code, 401)
	require.Equal(t, 1, len(c.Errors))
	assert.Equal(t, errors.ErrMissingApiKey.Error(), c.Errors[0].Error())
	assert.Equal(t, `{"message":"Unauthenticated","status":"error"}`, w.Body.String())
}

func TestRequireAPIkey_Wrong(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// test request, must instantiate a request first
	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header), // if you need to test headers
	}
	req.Header.Add("Authorization", "123456")
	c.Request = req

	test := RequireAPIkey("abcdefg")
	test(c)

	assert.Equal(t, w.Code, 401)
	require.Equal(t, 1, len(c.Errors))
	assert.Equal(t, errors.ErrMissingApiKey.Error(), c.Errors[0].Error())
	assert.Equal(t, `{"message":"Unauthenticated","status":"error"}`, w.Body.String())
}

func TestRequireAPIkey_NotSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// test request, must instantiate a request first
	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header), // if you need to test headers
	}
	c.Request = req

	test := RequireAPIkey("abcdefg")
	test(c)

	assert.Equal(t, w.Code, 401)
	require.Equal(t, 1, len(c.Errors))
	assert.Equal(t, errors.ErrMissingApiKey.Error(), c.Errors[0].Error())
	assert.Equal(t, `{"message":"Unauthenticated","status":"error"}`, w.Body.String())
}

func TestRequireCalendarSecret(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	url, err := url.Parse("https://localhost/test?token=calsecret-123456")
	require.NoError(t, err)

	req := &http.Request{
		URL:    url,
		Header: make(http.Header), // if you need to test headers
	}
	c.Request = req

	test := RequireCalendarSecret()
	test(c)

	value, ok := c.Get("token")
	require.True(t, ok)
	secret, ok := value.(string)
	require.True(t, ok)
	require.Equal(t, "calsecret-123456", secret)
}

func TestRequireCalendarSecret_Empty(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	url, err := url.Parse("https://localhost/test?token=")
	require.NoError(t, err)

	req := &http.Request{
		URL:    url,
		Header: make(http.Header), // if you need to test headers
	}
	c.Request = req

	test := RequireCalendarSecret()
	test(c)

	value, ok := c.Get("token")
	require.True(t, ok)
	secret, ok := value.(string)
	require.True(t, ok)
	require.Equal(t, "", secret)
}

func TestRequireCalendarSecret_NotSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	url, err := url.Parse("https://localhost/test")
	require.NoError(t, err)

	req := &http.Request{
		URL:    url,
		Header: make(http.Header), // if you need to test headers
	}
	c.Request = req

	test := RequireCalendarSecret()
	test(c)

	value, ok := c.Get("token")
	require.True(t, ok)
	secret, ok := value.(string)
	require.True(t, ok)
	require.Equal(t, "", secret)
}
