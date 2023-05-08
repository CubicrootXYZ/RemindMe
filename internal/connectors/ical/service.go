package ical

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/CubicrootXYZ/gologger"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

// Config holds the configuration for the service.
type Config struct {
	ICalDB   icaldb.Service
	Database database.Service

	BaseURL *url.URL // TODO set
}

type service struct {
	config *Config
	logger gologger.Logger

	stop chan bool
}

// New assembles a new ical connector service.
func New(config *Config, logger gologger.Logger) Service {
	return &service{
		config: config,
		logger: logger,
		stop:   make(chan bool, 1),
	}
}

func (service *service) Start() error {
	ticker := time.NewTicker(time.Minute * 15)
	for {
		service.refreshIcalInputs()

		select {
		case <-ticker.C:
			continue
		case <-service.stop:
			return nil
		}
	}
}

func (service *service) Stop() error {
	service.stop <- true
	service.logger.Infof("stopping")
	return nil
}

func (service *service) InputRemoved(inputType string, inputID uint) error {
	if inputType != InputType {
		return nil
	}

	err := service.config.ICalDB.DeleteIcalInput(inputID)
	if errors.Is(err, icaldb.ErrNotFound) {
		return nil
	}

	return err
}

func (service *service) OutputRemoved(outputType string, outputID uint) error {
	if outputType != OutputType {
		return nil
	}

	err := service.config.ICalDB.DeleteIcalOutput(outputID)
	if errors.Is(err, icaldb.ErrNotFound) {
		return nil
	}

	return err
}

func (service *service) NewOutput(channelID uint) (*icaldb.IcalOutput, string, error) {
	icalOutput, err := service.config.ICalDB.NewIcalOutput(&icaldb.IcalOutput{})
	if err != nil {
		return nil, "", err
	}

	err = service.config.Database.AddOutputToChannel(channelID, &database.Output{
		ChannelID:  channelID,
		OutputType: OutputType,
		OutputID:   icalOutput.ID,
		Enabled:    true,
	})
	if err != nil {
		return nil, "", err
	}

	icalURL := service.config.BaseURL.JoinPath(fmt.Sprintf("/ical/%d", icalOutput.ID))
	icalURL.RawQuery = url.Values{"token": []string{icalOutput.Token}}.Encode()

	return icalOutput, icalURL.String(), nil
}

func (service *service) GetOutput(outputID uint) (*icaldb.IcalOutput, string, error) {
	output, err := service.config.ICalDB.GetIcalOutputByID(outputID)
	if err != nil {
		if errors.Is(err, icaldb.ErrNotFound) {
			return nil, "", ErrNotFound
		}
		return nil, "", err
	}

	icalURL := service.config.BaseURL.JoinPath(fmt.Sprintf("/ical/%d", output.ID))
	icalURL.RawQuery = url.Values{"token": []string{output.Token}}.Encode()

	return output, icalURL.String(), nil
}

func (service *service) SendReminder(*daemon.Event, *daemon.Output) error {
	// Not supported.
	return nil
}
func (service *service) SendDailyReminder(*daemon.DailyReminder, *daemon.Output) error {
	// Not supported.
	return nil
}
