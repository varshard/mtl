package errors

import "net/http"

type Err interface {
	Error() string
	ErrCode() int
}

type ErrNotFound struct {
	error
}

func (err ErrNotFound) ErrCode() int {
	return http.StatusNotFound
}

func NewErrNotFound(err error) ErrNotFound {
	return ErrNotFound{err}
}

type ErrInvalidInput struct {
	error
}

func (err ErrInvalidInput) ErrorCode() int {
	return http.StatusBadRequest
}

func NewErrInvalidInput(err error) ErrInvalidInput {
	return ErrInvalidInput{err}
}

type ErrUnexpected struct {
	error
}

func (err ErrUnexpected) ErrorCode() int {
	return http.StatusInternalServerError
}

func NewErrUnexpected(err error) ErrUnexpected {
	return ErrUnexpected{err}
}
