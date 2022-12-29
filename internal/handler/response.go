package handler

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
)

// ResponseMessage is a data type for response messages send out via the api.
type ResponseMessage string

const (
	ResponseMessageInternalServerError = ResponseMessage("sorry, that went wrong on the server side")
	ResponseMessageNotFound            = ResponseMessage("entity not found")
	ResponseMessageNoID                = ResponseMessage("missing ID in request")
	ResponseMessageUnauthorized        = ResponseMessage("Unauthorized")
	ResponseMessageUnknownType         = ResponseMessage("type is not known")
	ResponseMessageMissingURL          = ResponseMessage("url is missing")
	ResponseMessageInvalidData         = ResponseMessage("invalid data")
)

type calendarResponse struct {
	ID                uint   `json:"id"`         // Internal id
	UserIdentifier    string `json:"user_id"`    // Matrix user identifier
	Token             string `json:"token"`      // Secret token to get the calendar file
	ChannelIdentifier string `json:"channel_id"` // Matrix channel identifier
}

type channelResponse struct {
	ID                uint        `json:"id"` // Internal id
	Created           time.Time   `json:"created"`
	ChannelIdentifier string      `json:"channel_id"` // Matrix channel identifier
	UserIdentifier    string      `json:"user_id"`    // Matrix user identifier
	TimeZone          string      `json:"timezone" default:""`
	DailyReminder     bool        `json:"daily_reminder"` // Whether the daily reminder is activated or not
	Role              *roles.Role `json:"role" enums:"user,admin" extensions:"x-nullable"`
}

type userResponse struct {
	UserIdentifier string            `json:"user_id"` // Matrix user identifier
	Blocked        bool              `json:"blocked"`
	Channels       []channelResponse `json:"channels"` // All channels known with the user
	Comment        string            `json:"comment"`
}

type thirdPartyResourceResponse struct {
	ID          uint   `json:"id"`   // Internal id
	Type        string `json:"type"` // The resources type
	ResourceURL string `json:"url"`  // The resources URL
}
