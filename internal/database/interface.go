package database

import (
	"errors"
	"time"

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

//go:generate mockgen -destination=service_mock.go -package=database . Service

// Service defines a database service interface.
type Service interface {
	// Channel
	NewChannel(*Channel) (*Channel, error)

	GetChannelByID(uint) (*Channel, error)

	AddInputToChannel(uint, *Input) error
	AddOutputToChannel(uint, *Output) error

	RemoveInputFromChannel(channelID, inputID uint) error
	RemoveOutputFromChannel(channelID, outputID uint) error

	UpdateChannel(channel *Channel) (*Channel, error)

	// Input
	GetInputByID(uint) (*Input, error)

	// Output
	GetOutputByID(uint) (*Output, error)

	// Event
	NewEvent(*Event) (*Event, error)

	GetEventsByChannel(uint) ([]Event, error)
	GetEventsPending() ([]Event, error)

	UpdateEvent(*Event) (*Event, error)
}

// Channel is the centerpiece orchestrating in- and outputs.
type Channel struct {
	gorm.Model
	Description       string
	DailyReminder     *uint // minutes from midnight when to send the daily reminder. Null to deactivate.
	Inputs            []Input
	Outputs           []Output
	LastDailyReminder *time.Time
}

// Input takes in data.
type Input struct {
	gorm.Model
	ChannelID uint
	Channel   Channel
	InputType string
	InputID   uint
	Enabled   bool
}

// Output takes data and moves it elsewhere.
type Output struct {
	gorm.Model
	ChannelID  uint
	Channel    Channel
	OutputType string
	OutputID   uint
	Enabled    bool
}

// Event holds information about an event
type Event struct {
	gorm.Model
	Time           time.Time `gorm:"index"`
	Duration       time.Duration
	Message        string
	Active         bool `gorm:"index"`
	RepeatInterval *time.Duration
	RepeatUntil    *time.Time
	ChannelID      uint
	Channel        Channel
	InputID        *uint
	Input          *Input
}
