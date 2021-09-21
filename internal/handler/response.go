package handler

// ResponseMessage is a data type for response messages send out via the api.
type ResponseMessage string

const (
	ResponseMessageInternalServerError = ResponseMessage("sorry, that went wrong on the server side")
	ResponseMessageNotFound            = ResponseMessage("entity not found")
)

type calendarResponse struct {
	ID      uint
	User    string
	Channel string
}

func responseMessageInvalidParameter(parameter string) ResponseMessage {
	return ResponseMessage("parameter " + parameter + " is missing or of wrong format")
}
