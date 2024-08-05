package httpe

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StatusError interface {
	json.Marshaler
	StatusCode() int
	GetError() error
	Unwrap() error
	error
}

type statusError struct {
	Err    error
	Status int
}

func (se statusError) GetError() error {
	return se.Err
}

func (se statusError) MarshalJSON() ([]byte, error) {
	return json.Marshal(se.Err.Error())
}

func (e statusError) Unwrap() error {
	return e.Err
}

func (e statusError) StatusCode() int {
	if e.Status == 0 {
		return http.StatusInternalServerError
	}
	return e.Status
}

func (e statusError) Error() string {
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
		httpError: statusError{
			Err:    base,
			Status: status,
		},
	}
}
