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
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/format"
)

func (service *service) refreshIcalInputs() {
	// TODO test
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

		_, err = service.config.ICalDB.UpdateIcalInput(&input)
		if err != nil {
			l.Infof("failed updating input in database: %v", err)
			continue
		}
	}
}

func (service *service) refreshIcalInput(input *database.IcalInput) error {
	// TODO test
	if input.Disabled {
		return nil
	}

	i, err := service.config.Database.GetInputByType(input.ID, InputType)
	if err != nil {
		return fmt.Errorf("can not get input for iCal input: %w", err)
	}

	content, err := getFileContent(input.URL)
	if err != nil {
		return fmt.Errorf("can not fetch resource: %w", err)
	}

	events, err := format.EventsFromIcal(content, &format.EventOpts{
		// TODO make configurable
		EventDelay:      time.Duration(0),
		DefaultDuration: time.Minute * 5,

		InputID:   i.ID,
		ChannelID: i.ChannelID,
	})
	if err != nil {
		return fmt.Errorf("failed to load events from iCal string: %w", err)
	}

	err = service.config.Database.NewEvents(events)
	if err != nil {
		return fmt.Errorf("failed to insert events to database: %w", err)
	}

	return nil
}

func getFileContent(url string) (string, error) {
	// TODO test
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
