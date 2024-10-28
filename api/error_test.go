package vaku

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// errInject is used when injecting errors in tests.
	errInject = errors.New("injected error")
)

func TestNewWrapErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		giveMsg  string
		giveIs   error
		giveWrap error
		wantErr  *wrapErr
	}{
		{
			name:     "nil all",
			giveMsg:  "",
			giveIs:   nil,
			giveWrap: nil,
			wantErr: &wrapErr{
				msg:   "unknown error",
				is:    ErrUnknownError,
				wraps: nil,
			},
		},
		{
			name:     "nil msg and is",
			giveMsg:  "",
			giveIs:   nil,
			giveWrap: errInject,
			wantErr: &wrapErr{
				msg:   fmt.Sprintf("%v: %v", ErrUnknownError, errInject),
				is:    ErrUnknownError,
				wraps: errInject,
			},
		},
		{
			name:     "nil is and wrap",
			giveMsg:  "random error",
			giveIs:   nil,
			giveWrap: nil,
			wantErr: &wrapErr{
				msg:   "random error",
				is:    errors.New("random error"),
				wraps: nil,
			},
		},
		{
			name:     "nil msg and wrap",
			giveMsg:  "",
			giveIs:   errInject,
			giveWrap: nil,
			wantErr: &wrapErr{
				msg:   errInject.Error(),
				is:    errInject,
				wraps: nil,
			},
		},
		{
			name:     "msg and nil wrap",
			giveMsg:  "random error",
			giveIs:   errInject,
			giveWrap: nil,
			wantErr: &wrapErr{
				msg:   fmt.Sprintf("%v: %v", "random error", errInject),
				is:    errInject,
				wraps: nil,
			},
		},
		{
			name:     "standard error",
			giveMsg:  "context here",
			giveIs:   errors.New("standard error"),
			giveWrap: errInject,
			wantErr: &wrapErr{
				msg:   fmt.Sprintf("%v: %v: %v", "context here", errors.New("standard error"), errInject),
				is:    errors.New("standard error"),
				wraps: errInject,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := newWrapErr(tt.giveMsg, tt.giveIs, tt.giveWrap)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCtxErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give error
		want []error
	}{
		{
			name: "nil error",
			give: nil,
			want: nil,
		},
		{
			name: "error",
			give: errInject,
			want: []error{ErrContext, errInject},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			compareErrors(t, ctxErr(tt.give), tt.want)
		})
	}
}
