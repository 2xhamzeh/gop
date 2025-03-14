package domain

import (
	"errors"
	"fmt"
)

// Code represents a domain error code.
// Can be used as an error type.
type ErrorCode string

const (
	CONFLICT_ERROR     = ErrorCode("conflict")
	INTERNAL_ERROR     = ErrorCode("internal")
	INVALID_ERROR      = ErrorCode("invalid")
	NOTFOUND_ERROR     = ErrorCode("not_found")
	UNAUTHORIZED_ERROR = ErrorCode("unauthorized")
	FORBIDDEN_ERROR    = ErrorCode("forbidden")
)

type Error struct {
	Code    ErrorCode
	Message string
	Fields  map[string][]string
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Errorf creates a domain error.
// It uses the fmt.Sprintf function to format the message.
func Errorf(code ErrorCode, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap uses errors.Join to combine the domain error with another error.
func (e *Error) Wrap(err error) error {
	return errors.Join(e, err)
}

