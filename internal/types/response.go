package types

// MessageErrorResponse is a default response
type MessageErrorResponse struct {
	Message string `example:"Unauthenticated"`
	Status  string `example:"error"`
} // @name ErrorResponse

// DataResponse is the default response for data
type DataResponse struct {
	Status string `example:"success"`
	Data   interface{}
}
