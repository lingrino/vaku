package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		giveE1  error
		giveE2  error
		wantErr error
	}{
		{
			name:    "nil",
			giveE1:  nil,
			giveE2:  nil,
			wantErr: nil,
		},
		{
			name:    "e1 nil",
			giveE1:  nil,
			giveE2:  errors.New("bar"),
			wantErr: errors.New("bar"),
		},
		{
			name:    "e2 nil",
			giveE1:  errors.New("foo"),
			giveE2:  nil,
			wantErr: errors.New("foo"),
		},
		{
			name:    "both",
			giveE1:  errors.New("foo"),
			giveE2:  errors.New("bar"),
			wantErr: errors.New("foo\nbar"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cli, _, _ := newTestCLI(t, nil)

			err := cli.combineErr(tt.giveE1, tt.giveE2)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		give       any
		giveErr    error
		giveFormat string
		wantOut    string
		wantErr    string
	}{
		{
			name: "nil",
			give: nil,
		},
		{
			name:       "text string",
			give:       "foo",
			giveFormat: "text",
			wantOut:    "foo\n",
		},
		{
			name:       "text list",
			give:       []string{"foo", "bar"},
			giveFormat: "text",
			wantOut:    "bar\nfoo\n",
		},
		{
			name: "text map",
			give: map[string]any{
				"foo": "fooValue",
				"bar": 100,
			},
			giveFormat: "text",
			wantOut:    "bar: 100\nfoo: fooValue\n",
		},
		{
			name: "text nested map",
			give: map[string]map[string]any{
				"foo": {
					"infoo": "fooValue",
					"inbar": 100,
				},
				"bar": {
					"hello": "world",
				},
			},
			giveFormat: "text",
			wantOut:    "bar\nhello: world\nfoo\ninbar: 100\ninfoo: fooValue\n",
		},
		{
			name:       "json string",
			give:       "foo",
			giveFormat: "json",
			wantOut:    "\"foo\"\n",
		},
		{
			name:       "json list",
			give:       []string{"foo", "bar"},
			giveFormat: "json",
			wantOut:    "[\n\"foo\",\n\"bar\"\n]\n",
		},
		{
			name: "json map",
			give: map[string]any{
				"foo": "fooValue",
				"bar": 100,
			},
			giveFormat: "json",
			wantOut:    "{\n\"bar\": 100,\n\"foo\": \"fooValue\"\n}\n",
		},
		{
			name: "json nested map",
			give: map[string]map[string]any{
				"foo": {
					"infoo": "fooValue",
					"inbar": 100,
				},
				"bar": {
					"hello": "world",
				},
			},
			giveFormat: "json",
			wantOut:    "{\n\"bar\": {\n\"hello\": \"world\"\n},\n\"foo\": {\n\"inbar\": 100,\n\"infoo\": \"fooValue\"\n}\n}\n",
		},
		{
			name:       "error text",
			give:       errors.New("test error"),
			giveFormat: "text",
			wantErr:    "ERROR: test error\n",
		},
		{
			name:       "error json",
			give:       errors.New("test error"),
			giveFormat: "json",
			wantErr:    "{\n\"error\": \"test error\"\n}\n",
		},
		{
			name:       "bad format",
			give:       "",
			giveFormat: "invalid",
			wantErr:    "ERROR: " + errOutputFormat.Error() + "\n",
		},
		{
			name:       "bad type",
			give:       5,
			giveFormat: "text",
			wantErr:    "ERROR: " + errOutputType.Error() + "\n",
		},
		{
			name:       "bad json",
			give:       func() {},
			giveFormat: "json",
			wantErr:    "ERROR: " + errJSONMarshal.Error() + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cli, outW, errW := newTestCLI(t, nil)
			cli.flagFormat = tt.giveFormat

			cli.output(tt.give)
			assert.Equal(t, tt.wantOut, outW.String())
			assert.Equal(t, tt.wantErr, errW.String())
		})
	}
}
