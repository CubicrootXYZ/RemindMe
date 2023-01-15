package database

import (
	"errors"

	"gorm.io/gorm"
)

// Errors returned by the service. List might not be complete.
var (
	ErrNotFound      = errors.New("entity not found")
	ErrInvalidConfig = errors.New("invalid config")
	ErrRolledBack    = errors.New("rolled back")
	ErrUnknownInput  = errors.New("unknown input type")
	ErrUnknownOutput = errors.New("unknown output type")
)

// Service defines a database service interface.
type Service interface {
	// Channel
	NewChannel(*Channel) (*Channel, error)

	GetChannelByID(uint) (*Channel, error)

	AddInputToChannel(uint, *Input) error
	AddOutputToChannel(uint, *Output) error

	RemoveInputFromChannel(channelID, inputID uint) error
	RemoveOutputFromChannel(channelID, outputID uint) error

	// Input
	GetInputByID(uint) (*Input, error)

	// Output
	GetOutputByID(uint) (*Output, error)
}

// Channel is the centerpiece orchestrating in- and outputs.
type Channel struct {
	gorm.Model
	Description   string
	DailyReminder *uint // minutes from midnight when to send the daily reminder. Null to deactivate.
	Inputs        []Input
	Outputs       []Output
}

// Input takes in data.
type Input struct {
	gorm.Model
	ChannelID uint
	Channel   *Channel
	InputType string
	InputID   uint
	Enabled   bool
}

// Output takes data and moves it elsewhere.
type Output struct {
	gorm.Model
	ChannelID  uint
	Channel    *Channel
	OutputType string
	OutputID   uint
	Enabled    bool
}
