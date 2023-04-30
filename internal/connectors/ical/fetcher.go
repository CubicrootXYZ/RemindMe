package ical

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
)

func (service *service) refreshIcalInputs() {
	service.logger.Debugf("refreshing iCal inputs now ...")

	f := false
	inputs, err := service.config.ICalDB.ListIcalInputs(&database.ListIcalInputsOpts{
		Disabled: &f,
	})
	if err != nil {
		service.logger.Err(err)
		return
	}

	for _, input := range inputs {
		input := input
		l := service.logger.WithField("iCal input ID", input.ID)

		now := time.Now()
		err := service.refreshIcalInput(&input)
		if err != nil {
			l.Infof("failed refreshing: %v", input.ID, err)
			if input.LastRefresh != nil && time.Since(*input.LastRefresh) > time.Hour*48 {
				l.Infof("disabling input, no success refrehsing since 48 hours")
				input.Disabled = true
			}
		} else {
			input.LastRefresh = &now
			input.Disabled = false
		}

		// TODO
		//service.config.ICalDB.
	}
}

func (service *service) refreshIcalInput(input *database.IcalInput) error {
	if input.Disabled {
		return nil
	}

	_, err := getFileContent(input.URL)
	if err != nil {
		return fmt.Errorf("can not fetch resource: %w", err)
	}

	// TODO

	return nil
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
