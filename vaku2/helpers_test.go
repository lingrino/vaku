package vaku2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyJoin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give []string
		want string
	}{
		{
			give: []string{"/"},
			want: "/",
		},
		{
			give: []string{"a/"},
			want: "a/",
		},
		{
			give: []string{"b", ""},
			want: "b",
		},
		{
			give: []string{"a/b", "c"},
			want: "a/b/c",
		},
		{
			give: []string{"d/e/", "/f"},
			want: "d/e/f",
		},
		{
			give: []string{"/g/h/", "/i/"},
			want: "g/h/i/",
		},
		{
			give: []string{"/j/", "/k/l", "m"},
			want: "j/k/l/m",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, KeyJoin(tt.give...))
		})
	}
}

func TestPathJoin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give []string
		want string
	}{
		{
			give: []string{"/"},
			want: "",
		},
		{
			give: []string{"a/"},
			want: "a",
		},
		{
			give: []string{"b", ""},
			want: "b",
		},
		{
			give: []string{"a/b", "c"},
			want: "a/b/c",
		},
		{
			give: []string{"d/e/", "/f"},
			want: "d/e/f",
		},
		{
			give: []string{"/g/h/", "/i/"},
			want: "g/h/i",
		},
		{
			give: []string{"/j/", "/k/l", "m"},
			want: "j/k/l/m",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, PathJoin(tt.give...))
		})
	}
}
