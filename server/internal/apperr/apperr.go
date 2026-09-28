package apperr

import "net/http"

type Error struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return e.Message }

var (
	ErrUnauthorized = &Error{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "authentication required"}
	ErrForbidden    = &Error{Status: http.StatusForbidden, Code: "FORBIDDEN", Message: "permission denied"}
	ErrNotFound     = &Error{Status: http.StatusNotFound, Code: "NOT_FOUND", Message: "resource not found"}
	ErrTooMany      = &Error{Status: http.StatusTooManyRequests, Code: "RATE_LIMITED", Message: "too many requests"}
	ErrTooLarge     = &Error{Status: http.StatusRequestEntityTooLarge, Code: "BODY_TOO_LARGE", Message: "request body too large"}
)

func BadRequest(msg string) *Error {
	return &Error{Status: http.StatusBadRequest, Code: "BAD_REQUEST", Message: msg}
}

func Conflict(msg string) *Error {
	return &Error{Status: http.StatusConflict, Code: "CONFLICT", Message: msg}
}
