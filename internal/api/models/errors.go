package models

import "fmt"

// ErrMissingRequiredField is an error type that represents a missing required field in a struct.
type ErrMissingRequiredField string

// Error() returns the error message for the ErrMissingRequiredField type.
func (e ErrMissingRequiredField) Error() string {
	return fmt.Sprintf("missing required field: %s", string(e))
}
