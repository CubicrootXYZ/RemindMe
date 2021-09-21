package errors

import "errors"

// Matrix related errors
var ErrMatrixClientNotInitialized = errors.New("matrix client is not initialized, can not perform action")
var ErrMatrixEventWrongType = errors.New("the received event is of wrong type")

// Generic errors
var ErrIdNotSet = errors.New("ID is not set")

// Api errors
var ErrAPIkeyCriteriaNotMet = errors.New("the api key does not met the minimum criteria (> 20 signs)")

// Gin errors
var ErrMissingApiKey = errors.New("unauthenticated")
var ErrMissingID = errors.New("Can not get ID from context")
