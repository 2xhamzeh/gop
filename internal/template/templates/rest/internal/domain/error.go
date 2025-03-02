package domain

import (
	"fmt"
)

const (
	CONFLICT_ERROR     = "conflict"
	INTERNAL_ERROR     = "internal"
	INVALID_ERROR      = "invalid"
	NOTFOUND_ERROR     = "not_found"
	UNAUTHORIZED_ERROR = "unauthorized"
)

type Error struct {
	Code    string
	Message string
	Fields  []string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

// Factory function for domain errors
func Errorf(code, format string, args ...any) error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Factory function for domain errors with fields
func ErrorfWithFields(code, format string, fields []string, args ...any) error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Fields:  fields,
	}
}
