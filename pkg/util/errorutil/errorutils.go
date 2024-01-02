package errorutil

import (
	"fmt"
	"runtime"
)

// CustomError struct for detailed error information
type CustomError struct {
	Type        ErrorType
	Msg         string
	OriginalErr error
	File        string
	Line        int
}

// Error implements the error interface
func (e CustomError) Error() string {
	return fmt.Sprintf("[%s] %s (File: %s, Line: %d)", e.Type, e.Msg, e.File, e.Line)
}

// Unwrap returns the original error (if available) for standard library error handling
func (e CustomError) Unwrap() error {
	return e.OriginalErr
}

// New creates a new CustomError with stack trace information
func New(errType ErrorType, msg string, originalErr error) error {
	_, file, line, _ := runtime.Caller(1)
	return CustomError{
		Type:        errType,
		Msg:         msg,
		OriginalErr: originalErr,
		File:        file,
		Line:        line,
	}
}

// IsType checks if the error is of a specific ErrorType
func IsType(err error, errType ErrorType) bool {
	var customErr CustomError
	if asErr := As(err, &customErr); asErr {
		return customErr.Type == errType
	}
	return false
}

// As attempts to cast an error to CustomError
func As(err error, target interface{}) bool {
	t, ok := target.(*CustomError)
	if !ok {
		return false
	}
	customErr, ok := err.(CustomError)
	if !ok {
		// Support for wrapped errors
		wrappedErr, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		return As(wrappedErr.Unwrap(), target)
	}
	*t = customErr
	return true
}

// Additional utility functions can be added here, like logging errors
