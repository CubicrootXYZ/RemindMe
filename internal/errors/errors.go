package errors

import "errors"

// Matrix related errors
var ErrMatrixClientNotInitialized = errors.New("matrix client is not initialized, can not perform action")
var ErrMatrixEventWrongType = errors.New("the received event is of wrong type")

// Generic errors
var ErrEmptyChannel = errors.New("the given channel is empty")
var ErrIdNotSet = errors.New("ID is not set")

// Api errors
var ErrAPIkeyCriteriaNotMet = errors.New("the api key does not met the minimum criteria (> 20 signs)")

// Gin errors
var ErrMissingApiKey = errors.New("unauthenticated")
var ErrMissingID = errors.New("can not get ID from context")
var ErrMissingIDString = errors.New("can not get ID (string) from context")

// Encryption errors
var ErrOlmNotSetUp = errors.New("not set up to handle encryption with olm")
var ErrRoomNotEncrypted = errors.New("the given room is not encrpyted")
