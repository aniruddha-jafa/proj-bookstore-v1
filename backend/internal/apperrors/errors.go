package apperrors

import "net/http"

// Errors
type HttpError struct {
	Status  int
	Message string
}

func NewHttpError(statusCode int, message string) HttpError {
	return HttpError{statusCode, message}
}

func (h HttpError) Error() string {
	return h.Message
}

var ErrNotFound = HttpError{Status: http.StatusNotFound, Message: "resource not found"}
var ErrInvalidParams = HttpError{Status: http.StatusBadGateway, Message: "invalid params"}
var ErrEmailTaken = HttpError{Status: http.StatusConflict, Message: "email is taken"}
var ErrInvalidCredentials = HttpError{Status: http.StatusUnauthorized, Message: "Invalid credentials"}
var ErrMissingAuthorizationHeader = HttpError{Status: http.StatusUnauthorized, Message: "Missing authorization header"}
var ErrMissingBearerToken = HttpError{Status: http.StatusUnauthorized, Message: "Missing bearer token"}
var ErrTokenRevoked = HttpError{Status: http.StatusUnauthorized, Message: "Token revoked"}
var ErrTokenExpired = HttpError{Status: http.StatusUnauthorized, Message: "Token expired"}
var ErrRefreshTokenNotFound = HttpError{Status: http.StatusUnauthorized, Message: "Refresh token not found"}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
