package task_client

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBadRequest    = errors.New("bad request")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrNotFound      = errors.New("not found")
	ErrConflict      = errors.New("conflict")
	ErrUnprocessable = errors.New("unprocessable entity")
	ErrServer        = errors.New("server error")
)

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func (e *APIError) Unwrap() error {
	switch e.StatusCode {
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusConflict:
		return ErrConflict
	case http.StatusUnprocessableEntity:
		return ErrUnprocessable
	default:
		if e.StatusCode >= http.StatusInternalServerError {
			return ErrServer
		}
		return nil
	}
}
