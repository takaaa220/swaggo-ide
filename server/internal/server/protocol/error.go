package protocol

import "fmt"

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Error returns a string representation of the error.
func (e *ResponseError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// NewResponseError creates a ResponseError with a specific code and message.
func NewResponseError(code LSPCodeError, message string, data any) *ResponseError {
	return &ResponseError{
		Code:    int(code),
		Message: message,
		Data:    data,
	}
}
