package validator

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

type validator struct {
	message string
	errors  map[string][]string
}

func New() *validator {
	return &validator{errors: make(map[string][]string)}
}

// Error is returned by e.Valid() if there are any errors.
type Error struct {
	Message string
	Fields  map[string][]string
}

func (e *Error) Error() string {
	return e.Message
}

// AddError adds an error message to the list of errors for a specific field.
func (v *validator) AddError(field, message string) {
	v.errors[field] = append(v.errors[field], message)
}

// Message sets a custom error message for the error returned by e.Valid().
// Only the last message set will be used.
func (v *validator) Message(message string) {
	v.message = message
}

// Validate returns the map of errors or nil if there are no errors.
func (v *validator) Valid() error {
	if len(v.errors) == 0 && v.message == "" {
		return nil
	}
	if v.message == "" {
		v.message = "invalid input"
	}
	return &Error{Message: v.message, Fields: v.errors}
}

// CheckField checks if a condition is met and adds an error message to the list of errors for a specific field if it's not.
func (v *validator) CheckField(ok bool, field, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

// Helper methods

func (v *validator) NotBlank(value, field, message string) {
	v.CheckField(strings.TrimSpace(value) != "", field, message)
}

func (v *validator) Email(value, field, message string) {
	_, err := mail.ParseAddress(value)
	v.CheckField(err == nil, field, message)
}

func (v *validator) MinRunes(value string, min int, field, message string) {
	v.CheckField(utf8.RuneCountInString(value) >= min, field, message)
}

func (v *validator) MaxRunes(value string, max int, field, message string) {
	v.CheckField(utf8.RuneCountInString(value) <= max, field, message)
}

func (v *validator) BetweenRunes(value string, min, max int, field, message string) {
	v.CheckField(utf8.RuneCountInString(value) >= min && utf8.RuneCountInString(value) <= max, field, message)
}
