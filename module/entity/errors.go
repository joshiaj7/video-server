package entity

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// General
	ErrorBadRequest = NewError("Bad Request", http.StatusBadRequest)

	ErrorParamType = NewError("Wrong param type", http.StatusUnprocessableEntity)

	ErrorFileNotFound    = NewError("File not found", http.StatusNotFound)
	ErrorFileExists      = NewError("File exists", http.StatusConflict)
	ErrorFileUnsupported = NewError("File unsupported", http.StatusUnsupportedMediaType)
)

type RequestError struct {
	StatusCode int
	Err        error
}

func (r RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func NewError(message string, code int) error {
	return RequestError{
		StatusCode: code,
		Err:        errors.New(message),
	}
}
