package response

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
