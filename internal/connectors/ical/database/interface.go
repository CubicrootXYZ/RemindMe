package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// List of common errors returned by the package.
var (
	ErrNotFound = errors.New("not found")
)

// IcalInput holds information about an iCal resource that can be fetched.
type IcalInput struct {
	gorm.Model
	URL         string
	LastRefresh *time.Time
	Disabled    bool // Disabled if fetching failed for to long.
}

// IcalOutput holds information about an iCal resource holding channel events.
type IcalOutput struct {
	gorm.Model
	Token string
}

//go:generate mockgen -destination=service_mock.go -package=database . Service

// Service provides a database service for the iCal connector.
type Service interface {
	NewIcalInput(*IcalInput) (*IcalInput, error)
	GetIcalInputByID(id uint) (*IcalInput, error)
	ListIcalInputs(*ListIcalInputsOpts) ([]IcalInput, error)
	UpdateIcalInput(entity *IcalInput) (*IcalInput, error)
	DeleteIcalInput(id uint) error

	NewIcalOutput(*IcalOutput) (*IcalOutput, error)
	GetIcalOutputByID(id uint) (*IcalOutput, error)
	GenerateNewToken(*IcalOutput) (*IcalOutput, error)
	DeleteIcalOutput(id uint) error
}

// ListIcalInputsOpts holds options for listing iCal inputs.
type ListIcalInputsOpts struct {
	Disabled *bool
}
