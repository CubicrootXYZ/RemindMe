package icalimporter

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	ical "github.com/arran4/golang-ical"
)

func getNameFromEvent(event *ical.VEvent) string {
	for _, property := range []ical.ComponentProperty{ical.ComponentPropertyDescription, ical.ComponentPropertySummary} {
		prop := event.GetProperty(property)
		if prop == nil {
			continue
		}

		return prop.Value
	}

	return ""
}

func getFileContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("bad status code: " + strconv.Itoa(resp.StatusCode))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
