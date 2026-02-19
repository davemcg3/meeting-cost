package errors

// Standard error codes used across the backend. These mirror the API error
// codes defined in the backend contracts so handlers can map DomainErrors
// directly to structured API responses.

const (
	// Generic codes
	CodeValidation   = "VALIDATION_ERROR"
	CodeNotFound     = "NOT_FOUND"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeConflict     = "CONFLICT"
	CodeInternal     = "INTERNAL_ERROR"
	CodeRateLimit    = "RATE_LIMIT_EXCEEDED"
	CodeBadRequest   = "BAD_REQUEST"

	// Domain-specific codes
	CodeMeetingActive   = "MEETING_ACTIVE"
	CodeMeetingNotFound = "MEETING_NOT_FOUND"
)

