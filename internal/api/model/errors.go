// Package model provides data structures and error types for the recipe generator application.
package model

import "fmt"

// ErrMissingRequiredField is an error type that represents a missing required field in a struct.
// It implements the error interface and can be used to indicate which specific field is missing.
type ErrMissingRequiredField string

// Error implements the error interface for ErrMissingRequiredField.
// It returns a formatted error message indicating which field is missing.
func (e ErrMissingRequiredField) Error() string {
	return fmt.Sprintf("missing required field: %s", string(e))
}
