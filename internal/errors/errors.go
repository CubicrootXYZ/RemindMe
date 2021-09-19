package errors

import "errors"

var ErrMatrixClientNotInitialized = errors.New("matrix client is not initialized, can not perform action")
var ErrMatrixEventWrongType = errors.New("the received event is of wrong type")
var ErrEmptyChannel = errors.New("the given channel is empty")
var ErrIdNotSet = errors.New("ID is not set")
