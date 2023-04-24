package ical

import (
	"errors"

	"github.com/CubicrootXYZ/gologger"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

// Config holds the configuration for the service.
type Config struct {
	ICalDB   icaldb.Service
	Database database.Service
}

type service struct {
	config *Config
	logger gologger.Logger
}

// New assembles a new ical connector service.
func New(config *Config, logger gologger.Logger) Service {
	return &service{
		config: config,
		logger: logger,
	}
}

func (service *service) Start() error {
	// TODO
	return nil
}

func (service *service) Stop() error {
	// TODO
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

func (service *service) NewOutput(channelID uint) (*icaldb.IcalOutput, error) {
	icalOutput, err := service.config.ICalDB.NewIcalOutput(&icaldb.IcalOutput{})
	if err != nil {
		return nil, err
	}

	err = service.config.Database.AddOutputToChannel(channelID, &database.Output{
		ChannelID:  channelID,
		OutputType: OutputType,
		OutputID:   icalOutput.ID,
		Enabled:    true,
	})
	if err != nil {
		return nil, err
	}

	return icalOutput, nil
}

func (service *service) GetOutput(outputID uint) (*icaldb.IcalOutput, error) {
	output, err := service.config.ICalDB.GetIcalOutputByID(outputID)
	if errors.Is(err, icaldb.ErrNotFound) {
		return nil, ErrNotFound
	}

	return output, err
}
