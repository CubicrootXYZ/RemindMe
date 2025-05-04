package ical

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

const (
	defaultEventDuration = time.Minute * 5
)

func (service *service) refreshIcalInputs() {
	service.logger.Debug("refreshing iCal inputs")

	f := false
	inputs, err := service.config.ICalDB.ListIcalInputs(&icaldb.ListIcalInputsOpts{
		Disabled: &f,
	})
	if err != nil {
		service.logger.Error("failed to list iCal inputs", "error", err)
		return
	}

	for _, input := range inputs {
		l := service.logger.With("ical.input_id", input.ID)

		if input.LastRefresh != nil && time.Since(*input.LastRefresh) < service.config.RefreshInterval {
			l.Debug("skipping input as has not yet reached the refresh interval")
			continue
		}

		now := time.Now()
		service.metricLastRefresh.
			WithLabelValues(strconv.FormatUint(uint64(input.ID), 10)).
			Set(float64(now.Unix()))

		err := service.refreshIcalInput(&input)
		if err != nil {
			service.metricErrorCount.
				WithLabelValues(strconv.FormatUint(uint64(input.ID), 10)).
				Inc()

			l.Info("failed refreshing input", "error", err)
			if input.LastRefresh != nil && time.Since(*input.LastRefresh) > time.Hour*48 {
				l.Info("disabling input, no successful refresh in 48 hours")
				input.Disabled = true
			}
		} else {
			input.LastRefresh = &now
			input.Disabled = false
		}

		_, err = service.config.ICalDB.UpdateIcalInput(&input)
		if err != nil {
			service.metricErrorCount.
				WithLabelValues(strconv.FormatUint(uint64(input.ID), 10)).
				Inc()

			l.Info("failed updating input in database", "error", err)
			continue
		}
	}
}

func (service *service) refreshIcalInput(input *icaldb.IcalInput) error {
	if input.Disabled {
		return nil
	}

	i, err := service.config.Database.GetInputByType(input.ID, InputType)
	if err != nil {
		return fmt.Errorf("can not get input for iCal input: %w", err)
	}

	content, err := getFileContent(input.URL, service.logger)
	if err != nil {
		return fmt.Errorf("can not fetch resource: %w", err)
	}

	events, err := format.EventsFromIcal(content, &format.EventOpts{
		// TODO make configurable
		EventDelay:      time.Duration(0),
		DefaultDuration: defaultEventDuration,

		InputID:   i.ID,
		ChannelID: i.ChannelID,
	})
	if err != nil {
		return fmt.Errorf("failed to load events from iCal string: %w", err)
	}

	existingEvents, err := service.config.Database.ListEvents(&database.ListEventsOpts{
		InputID: &i.ID,
	})
	if err != nil {
		return fmt.Errorf("can not list existing events: %w", err)
	}

	newEvents := make([]database.Event, 0)
	for _, event := range events {
		if event.ExternalReference == "" {
			continue
		}

		update := false
		for _, eEvent := range existingEvents {
			if eEvent.ExternalReference == "" {
				continue
			}

			if eEvent.ExternalReference == event.ExternalReference {
				update = true
				break
			}
		}

		if !update {
			newEvents = append(newEvents, event)
			continue
		}

		// TODO handle event updates.
	}

	err = service.config.Database.NewEvents(newEvents)
	if err != nil {
		return fmt.Errorf("failed to insert events to database: %w", err)
	}

	return nil
}

func getFileContent(url string, logger *slog.Logger) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/calendar")
	logger.Debug("making request", "url", url, "http.method", "GET")

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
