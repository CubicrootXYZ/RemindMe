package database

import "gorm.io/gorm"

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

type Service interface {
	NewIcalInput(*IcalInput) (*IcalInput, error)
	DeleteIcalInput(id uint) error

	NewIcalOutput(*IcalOutput) (*IcalOutput, error)
	GenerateNewToken(*IcalOutput) (*IcalOutput, error)
	DeleteIcalOutput(id uint) error
}
