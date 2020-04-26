package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		give       interface{}
		giveFormat string
		wantOut    string
		wantErr    string
	}{
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
			give: map[string]interface{}{
				"foo": "fooValue",
				"bar": 100,
			},
			giveFormat: "text",
			wantOut:    "bar => 100\nfoo => fooValue\n",
		},
		{
			name: "text nested map",
			give: map[string]map[string]interface{}{
				"foo": {
					"infoo": "fooValue",
					"inbar": 100,
				},
				"bar": {
					"hello": "world",
				},
			},
			giveFormat: "text",
			wantOut:    "bar\nhello => world\nfoo\ninbar => 100\ninfoo => fooValue\n",
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
			give: map[string]interface{}{
				"foo": "fooValue",
				"bar": 100,
			},
			giveFormat: "json",
			wantOut:    "{\n\"bar\": 100,\n\"foo\": \"fooValue\"\n}\n",
		},
		{
			name: "json nested map",
			give: map[string]map[string]interface{}{
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
			name:       "bad format",
			give:       "",
			giveFormat: "invalid",
			wantErr:    errOutputFormat + "\n",
		},
		{
			name:       "bad type",
			give:       5,
			giveFormat: "text",
			wantErr:    errOutputType + "\n",
		},
		{
			name:       "bad json",
			give:       func() {},
			giveFormat: "json",
			wantErr:    errJSONMarshal + "\n",
		},
	}

	for _, tt := range tests {
		tt := tt
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
