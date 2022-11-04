package icalimporter

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/calendar")

	client := &http.Client{}
	resp, err := client.Do(req)
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
