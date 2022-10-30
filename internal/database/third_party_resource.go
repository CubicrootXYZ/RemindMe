package database

import (
	"strings"

	"gorm.io/gorm"
)

type ThirdPartyResourceType string

func (resourceType ThirdPartyResourceType) String() string {
	return string(resourceType)
}

// ThirdPartyResourceTypeFromString parses the resource type from a string
func ThirdPartyResourceTypeFromString(resourceType string) (ThirdPartyResourceType, error) {
	switch strings.ToLower(strings.TrimSpace(resourceType)) {
	case "ical":
		return ThirdPartyResourceTypeIcal, nil
	}

	return ThirdPartyResourceType(""), ErrThirdPartyResourceTypeUnknown
}

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

// GetThirdPartyResources lists all resources with the given type
func (d *Database) GetThirdPartyResources(resourceType ThirdPartyResourceType) ([]ThirdPartyResource, error) {
	resources := make([]ThirdPartyResource, 0)
	err := d.db.Find(&resources, "type = ?", resourceType.String()).Error

	return resources, err
}
