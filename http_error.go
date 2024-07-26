package httpe

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StatusError struct {
	Err    error
	Status int
}

func (se StatusError) MarshalJSON() ([]byte, error) {
	return json.Marshal(se.Err.Error())
}

func (e StatusError) Unwrap() error {
	return e.Err
}

func (e StatusError) StatusCode() int {
	if e.Status == 0 {
		return http.StatusInternalServerError
	}
	return e.Status
}

func (e StatusError) Error() string {
	return e.Err.Error()
}

type Error struct {
	baseError error
	httpError StatusError
}

func (e Error) Error() string {
	return e.baseError.Error()
}

func (e Error) Unwrap() error {
	return fmt.Errorf("%w %w", e.baseError, e.httpError)
}

func NewError(base error, status int) Error {
	return Error{
		baseError: base,
		httpError: StatusError{
			Err:    base,
			Status: status,
		},
	}
}
