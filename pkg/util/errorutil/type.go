package errorutil

// ErrorType categorizes the error
type ErrorType string

// Predefined error types
const (
	NotFoundError   ErrorType = "NotFoundError"
	ValidationError ErrorType = "ValidationError"
	DatabaseError   ErrorType = "DatabaseError"
	NetworkError    ErrorType = "NetworkError"
	// More error types can be added here
)
