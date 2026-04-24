package goaviatrix

import "fmt"

type StatusError struct {
	code int
	err  error
}

func (s StatusError) StatusCode() int {
	return s.code
}

func (s StatusError) Error() string {
	return s.err.Error()
}

func (s StatusError) Unwrap() error {
	return s.err
}

func NewStatusError(statusCode int, err error) StatusError {
	return StatusError{
		code: statusCode,
		err:  err,
	}
}

func NewStatusErrorf(statusCode int, format string, args ...any) StatusError {
	return StatusError{
		code: statusCode,
		err:  fmt.Errorf(format, args...),
	}
}
