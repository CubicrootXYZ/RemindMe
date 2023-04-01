package database

import (
	"errors"

	"gorm.io/gorm"
)

// List of common errors returned by the package.
var (
	ErrNotFound = errors.New("not found")
)

// IcalInput holds information about an iCal resource that can be fetched.
type IcalInput struct {
	gorm.Model
	URL string
}

// IcalOutput holds information about an iCal resource holding channel events.
type IcalOutput struct {
	gorm.Model
	Token string
}

// Service provides a database service for the iCal connector.
type Service interface {
	NewIcalInput(*IcalInput) (*IcalInput, error)
	GetIcalInputByID(id uint) (*IcalInput, error)
	DeleteIcalInput(id uint) error

	NewIcalOutput(*IcalOutput) (*IcalOutput, error)
	GenerateNewToken(*IcalOutput) (*IcalOutput, error)
	DeleteIcalOutput(id uint) error
}
