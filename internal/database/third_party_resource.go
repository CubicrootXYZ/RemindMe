package database

import "gorm.io/gorm"

type ThirdPartyResourceType string

// List of available third party resource types
var (
	ThirdPartyResourceTypeIcal = ThirdPartyResourceType("ICAL")
)

type ThirdPartyResource struct {
	gorm.Model
	Type        ThirdPartyResourceType
	ChannelID   uint `gorm:"index"`
	Channel     Channel
	ResourceURL string
}
