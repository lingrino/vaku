package vaku2

import "errors"

// Here are errors that are not specific to one file/function
var (
	// ErrDecodeSecret secret data cannot be extracted from a vault secret.
	ErrDecodeSecret = errors.New("decode secret")
	// ErrNilData returns when passed data is nil.
	ErrNilData = errors.New("nil data")
)

// wrapErr is a struct that implements the error interface and provides Is() and Unwrap() methods
// that allow go 1.13+ error features. The fmt.Errorf function does something similar but does not
// provide and Is() function which means you cannot use sentinel errors with added context and also
// wrap the returned error. context - https://golang.org/src/fmt/errors.go
type wrapErr struct {
	msg   string
	is    error
	wraps error
}

// newWrapErr returns a wrapErr with the msg, is, and wraps set
func newWrapErr(msg string, is, wraps error) *wrapErr {
	return &wrapErr{
		msg:   msg,
		is:    is,
		wraps: wraps,
	}
}

// Is compares an error to wrapErr.is
func (e *wrapErr) Is(target error) bool {
	return target == e.is
}

// Error() returns the error message
func (e *wrapErr) Error() string {
	return e.msg
}

// Unwrap returns the wrapped error
func (e *wrapErr) Unwrap() error {
	return e.wraps
}
