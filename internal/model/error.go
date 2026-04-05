package model

import "net/http"

type ErrorResponse struct {
	Message string `json:"message"`
}

type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewErrNotFound(message string) *AppError {
	return &AppError{Status: http.StatusNotFound, Message: message}
}

func NewErrBadRequest(message string) *AppError {
	return &AppError{Status: http.StatusBadRequest, Message: message}
}

func NewErrUnauthorized() *AppError {
	return &AppError{Status: http.StatusNotFound, Message: "Unauthorized"}
}

func NewConflictError(message string) *AppError {
	return &AppError{Status: http.StatusConflict, Message: message}
}
