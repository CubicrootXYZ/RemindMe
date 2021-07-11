package errors

import "errors"

var MatrixClientNotInitialized = errors.New("Matrix client is not initialized, can not perform action")
var MatrixEventWrongType = errors.New("the received event is of wrong type")
