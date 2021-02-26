// Source: https://github.com/benbjohnson/wtf/blob/main/error.go
package render

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	EINTERNAL       = "internal"
	EINVALID        = "invalid"
	ENOTFOUND       = "not_found"
	ENOTIMPLEMENTED = "not_implemented"
	EUNAUTHORIZED   = "unauthorized"
)

var codes = map[string]int{
	EINVALID:        http.StatusBadRequest,
	ENOTFOUND:       http.StatusNotFound,
	ENOTIMPLEMENTED: http.StatusNotImplemented,
	EUNAUTHORIZED:   http.StatusUnauthorized,
	EINTERNAL:       http.StatusInternalServerError,
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	// Extract error code & message.
	code, message := ParseErrorCode(err), ParseErrorMessage(err)

	Status(r, ConvertErrorCodeToStatusCode(code))
	JSON(w, r, &ErrorResponse{Error: message})
}

// ErrorResponse represents a JSON structure for error output.
type ErrorResponse struct {
	Error string `json:"error"`
}

type ApiError struct {
	// Machine-readable error code.
	Code string

	// Human-readable error message.
	Message string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("api error: code=%s message=%s", e.Code, e.Message)
}

// ParseErrorCode unwraps an application error and returns its code.
// Non-application errors always return EINTERNAL.
func ParseErrorCode(err error) string {
	var e *ApiError

	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}

	return EINTERNAL
}

// ParseErrorMessage unwraps an application error and returns its message.
// Non-application errors always return "Internal error".
func ParseErrorMessage(err error) string {
	var e *ApiError

	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}

	return "Internal error."
}

// Errorf is a helper function to return an Error with a given code and formatted message.
func Errorf(code string, format string, args ...interface{}) *ApiError {
	return &ApiError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// ErrorStatusCode returns the associated HTTP status code for an internal error code.
func ConvertErrorCodeToStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}
