package errors

import (
	"fmt"

	"github.com/google/uuid"
)

// DomainError represents a structured domain-level error that can be
// converted into an HTTP response and logged with context.
type DomainError struct {
	Code    string                 // Stable machine-readable error code
	Message string                 // Human-readable message (safe for clients)
	Details map[string]interface{} // Optional structured details
	Cause   error                  // Wrapped underlying error (not exposed directly)
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WithDetails attaches structured details to the error.
func (e *DomainError) WithDetails(details map[string]interface{}) *DomainError {
	e.Details = details
	return e
}

// WithCause wraps an underlying cause error.
func (e *DomainError) WithCause(cause error) *DomainError {
	e.Cause = cause
	return e
}

// Predefined generic domain errors.
var (
	ErrNotFound     = &DomainError{Code: CodeNotFound, Message: "resource not found"}
	ErrUnauthorized = &DomainError{Code: CodeUnauthorized, Message: "unauthorized"}
	ErrForbidden    = &DomainError{Code: CodeForbidden, Message: "forbidden"}
	ErrValidation   = &DomainError{Code: CodeValidation, Message: "validation failed"}
	ErrConflict     = &DomainError{Code: CodeConflict, Message: "resource conflict"}
)

// Helper constructors for common domain-specific errors.

func ErrPersonNotFound(id uuid.UUID) *DomainError {
	return &DomainError{
		Code:    "PERSON_NOT_FOUND",
		Message: fmt.Sprintf("person with ID %s not found", id),
		Details: map[string]interface{}{"person_id": id},
	}
}

func ErrOrganizationNotFound(id uuid.UUID) *DomainError {
	return &DomainError{
		Code:    "ORGANIZATION_NOT_FOUND",
		Message: fmt.Sprintf("organization with ID %s not found", id),
		Details: map[string]interface{}{"organization_id": id},
	}
}

func ErrMeetingNotFound(id uuid.UUID) *DomainError {
	return &DomainError{
		Code:    "MEETING_NOT_FOUND",
		Message: fmt.Sprintf("meeting with ID %s not found", id),
		Details: map[string]interface{}{"meeting_id": id},
	}
}

