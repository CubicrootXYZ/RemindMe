package response

// MessageErrorResponse is a default response
type MessageErrorResponse struct {
	Message string `example:"Unauthenticated" json:"message"`
	Status  string `example:"error"           json:"status"`
} // @name ErrorResponse

// MessageSuccessResponse is a default response
type MessageSuccessResponse struct {
	Message string `example:"Inserted new reminder" json:"message"`
	Status  string `example:"success"               json:"status"`
} // @name SuccessResponse

// DataResponse is the default response for data
type DataResponse struct {
	Status string `example:"success" json:"status"`
	Data   any    `json:"data"`
}
