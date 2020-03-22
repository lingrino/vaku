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

func TestPrefixList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		giveList   []string
		givePrefix string
		want       []string
	}{
		{
			giveList:   []string{"a"},
			givePrefix: "b",
			want:       []string{"b/a"},
		},
		{
			giveList:   []string{"/c/d/e/"},
			givePrefix: "/f/",
			want:       []string{"f/c/d/e/"},
		},
		{
			giveList:   []string{"/g/"},
			givePrefix: "h",
			want:       []string{"h/g/"},
		},
		{
			giveList:   []string{"i/j"},
			givePrefix: "i",
			want:       []string{"i/i/j"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.givePrefix, func(t *testing.T) {
			t.Parallel()

			PrefixList(tt.giveList, tt.givePrefix)

			assert.Equal(t, tt.want, tt.giveList)
		})
	}
}

func TestTrimListPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		giveList   []string
		givePrefix string
		want       []string
	}{
		{
			giveList:   []string{"a"},
			givePrefix: "b",
			want:       []string{"a"},
		},
		{
			giveList:   []string{"/c/d/e/"},
			givePrefix: "/c/",
			want:       []string{"d/e/"},
		},
		{
			giveList:   []string{"f/g"},
			givePrefix: "f",
			want:       []string{"g"},
		},
		{
			giveList:   []string{"i/j"},
			givePrefix: "k",
			want:       []string{"i/j"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.givePrefix, func(t *testing.T) {
			t.Parallel()

			TrimListPrefix(tt.giveList, tt.givePrefix)

			assert.Equal(t, tt.want, tt.giveList)
		})
	}
}
