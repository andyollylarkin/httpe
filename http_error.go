package httpe

import "fmt"

type HttpError struct {
	err    error
	status int
}

func (e HttpError) Unwrap() error {
	return e.err
}

func (e HttpError) Status() int {
	return e.status
}

func (e HttpError) Error() string {
	return fmt.Sprintf("Error: %s: status - %d", e.err.Error(), e.status)
}

type Error struct {
	baseError error
	httpError HttpError
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
		httpError: HttpError{
			err:    base,
			status: status,
		},
	}
}
