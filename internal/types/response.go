package types

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
)

// MessageErrorResponse is a default response
type MessageErrorResponse struct {
	Message string `json:"message" example:"Unauthenticated"`
	Status  string `json:"status" example:"error"`
} // @name ErrorResponse

// MessageSuccessResponse is a default response
type MessageSuccessResponse struct {
	Message string `json:"message" example:"Inserted new reminder"`
	Status  string `json:"status" example:"success"`
} // @name SuccessResponse

// DataResponse is the default response for data
type DataResponse struct {
	Status string      `json:"status" example:"success"`
	Data   interface{} `json:"data"`
}

// ResponseMessage is a data type for response messages send out via the api.
type ResponseMessage string

const (
	ResponseMessageInternalServerError = ResponseMessage("sorry, that went wrong on the server side")
	ResponseMessageNotFound            = ResponseMessage("entity not found")
	ResponseMessageNoID                = ResponseMessage("missing ID in request")
	ResponseMessageUnauthorized        = ResponseMessage("Unauthorized")
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
	Role              *roles.Role `json:"role" enums:"user,admin"`
}
