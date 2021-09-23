package types

// MessageErrorResponse is a default response
type MessageErrorResponse struct {
	Message string `example:"Unauthenticated"`
	Status  string `example:"error"`
} // @name ErrorResponse

// MessageSuccessResponse is a default response
type MessageSuccessResponse struct {
	Message string `example:"Inserted new reminder"`
	Status  string `example:"success"`
} // @name SuccessResponse

// DataResponse is the default response for data
type DataResponse struct {
	Status string `example:"success"`
	Data   interface{}
}
