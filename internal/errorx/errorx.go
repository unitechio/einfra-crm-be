
package errorx

import (
	"fmt"
	"runtime"

	"github.com/google/uuid"
)

// Error represents a custom error with additional context.
type Error struct {
	Code    int
	Message string
	TraceID string
	Stack   string
}

// Error returns the string representation of the error.
func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s, traceID: %s", e.Code, e.Message, e.TraceID)
}

// New creates a new Error with a unique trace ID.
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		TraceID: uuid.New().String(),
	}
}

// WithStack adds a stack trace to the error.
func (e *Error) WithStack() *Error {
	buf := make([]byte, 1<<16)
	length := runtime.Stack(buf, false)
	e.Stack = string(buf[:length])
	return e
}
