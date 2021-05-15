package vaku

import (
	"errors"
	"fmt"
)

// Errors that are not specific to one file/function.
var (
	// ErrContext when ctx.Err() is not nil.
	ErrContext = errors.New("context")
	// ErrDecodeSecret when secret data cannot be extracted from a vault secret.
	ErrDecodeSecret = errors.New("decode secret")
	// ErrJSONMarshal when secret data cannot be marshaled into json.
	ErrJSONMarshal = errors.New("json marshal")
	// ErrNilData when passed data is nil.
	ErrNilData = errors.New("nil data")
	// ErrUnknownError when returning an error with no data.
	ErrUnknownError = errors.New("unknown error")
)

// wrapErr is a struct that implements the error interface and provides Is() and Unwrap() methods
// that allow go 1.13+ error features. The fmt.Errorf function does something similar but does not
// provide an Is() function which means you cannot use sentinel errors with added context and also
// wrap the returned error. Context - https://golang.org/src/fmt/errors.go.
type wrapErr struct {
	msg   string
	is    error
	wraps error
}

// verify compliance with error interface.
var _ error = (*wrapErr)(nil)

// newWrapErr returns a wrapErr that merges defaults and input.
func newWrapErr(msg string, is, wraps error) *wrapErr {
	switch {
	case msg == "" && is == nil:
		is = ErrUnknownError
	case is == nil:
		is = errors.New(msg)
	}

	switch {
	case msg == "" && wraps == nil || msg == is.Error():
		msg = is.Error()
	case msg == "":
		msg = fmt.Sprintf("%v: %v", is, wraps)
	case wraps == nil:
		msg = fmt.Sprintf("%v: %v", msg, is)
	default:
		msg = fmt.Sprintf("%v: %v: %v", msg, is, wraps)
	}

	return &wrapErr{
		msg:   msg,
		is:    is,
		wraps: wraps,
	}
}

// Is compares an error to e.is.
func (e *wrapErr) Is(target error) bool {
	return target == e.is //nolint:errorlint
}

// Error() returns the error message.
func (e *wrapErr) Error() string {
	return e.msg
}

// Unwrap returns the wrapped error.
func (e *wrapErr) Unwrap() error {
	return e.wraps
}

// ctxErr returns a wrapped ErrContext if err != nil.
func ctxErr(err error) error {
	if err == nil {
		return nil
	}

	return newWrapErr("", ErrContext, err)
}
