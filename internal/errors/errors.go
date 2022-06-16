package errors

import "errors"

// Errors
var (
	// Matrix related errors
	ErrMatrixClientNotInitialized = errors.New("matrix client is not initialized, can not perform action")
	ErrMatrixEventWrongType       = errors.New("the received event is of wrong type")

	// Generic errors
	ErrEmptyChannel = errors.New("the given channel is empty")
	ErrIDNotSet     = errors.New("ID is not set")

	// Api errors
	ErrAPIkeyCriteriaNotMet = errors.New("the api key does not met the minimum criteria (> 20 signs)")

	// Gin errors
	ErrMissingAPIKey   = errors.New("unauthenticated")
	ErrMissingID       = errors.New("can not get ID from context")
	ErrMissingIDString = errors.New("can not get ID (string) from context")

	// Encryption errors
	ErrOlmNotSetUp      = errors.New("not set up to handle encryption with olm")
	ErrRoomNotEncrypted = errors.New("the given room is not encrpyted")

	// Config related errors
	ErrReactionsDisabled = errors.New("reactions are not enabled")
)
