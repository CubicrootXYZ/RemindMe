package icalimporter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ical "github.com/arran4/golang-ical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNameFromEvent_Description(t *testing.T) {
	event := &ical.VEvent{}
	event.SetProperty(ical.ComponentPropertyDescription, "eventname")

	assert.Equal(t, "eventname", getNameFromEvent(event))
}

func TestGetNameFromEvent_Summary(t *testing.T) {
	event := &ical.VEvent{}
	event.SetProperty(ical.ComponentPropertySummary, "eventname")

	assert.Equal(t, "eventname", getNameFromEvent(event))
}

func TestGetNameFromEvent_Empty(t *testing.T) {
	event := &ical.VEvent{}

	assert.Equal(t, "", getNameFromEvent(event))
}

func TestGetFileContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))

	content, err := getFileContent(server.URL)
	require.NoError(t, err)
	assert.Equal(t, "ok", content)
}

func TestGetFileContent_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	_, err := getFileContent(server.URL)
	assert.Error(t, err)
}
