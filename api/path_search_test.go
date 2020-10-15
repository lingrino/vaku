package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathSearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give        string
		giveSearch  string
		wantSuccess bool
		wantErr     []error
	}{
		{
			give:        "0/1",
			giveSearch:  "2",
			wantSuccess: true,
			wantErr:     nil,
		},
		{
			give:        "0/4/5",
			giveSearch:  "7",
			wantSuccess: true,
			wantErr:     nil,
		},
		{
			give:        "0/4/8",
			giveSearch:  "13",
			wantSuccess: false,
			wantErr:     nil,
		},
		{
			give:        "0/4/13/17",
			giveSearch:  "9",
			wantSuccess: true,
			wantErr:     nil,
		},
		{
			give:        "fake",
			giveSearch:  "searchstring",
			wantSuccess: false,
			wantErr:     nil,
		},
		{
			give:        "fakeempty",
			giveSearch:  "",
			wantSuccess: false,
			wantErr:     nil,
		},
		{
			give:        mountless,
			giveSearch:  "searchstring",
			wantSuccess: false,
			wantErr:     []error{ErrPathSearch, ErrPathRead, ErrRewritePath, ErrMountInfo, ErrNoMount},
		},
		{
			give:        "error/read/inject",
			giveSearch:  "searchstring",
			wantSuccess: false,
			wantErr:     []error{ErrPathSearch, ErrPathRead, ErrVaultRead},
		},
		{
			give:        "funcdata/read/inject",
			giveSearch:  "searchstring",
			wantSuccess: false,
			wantErr:     []error{ErrPathSearch, ErrJSONMarshal},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				prefix := prefix
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					success, err := sharedVaku.PathSearch(PathJoin(prefix, tt.give), tt.giveSearch)

					compareErrors(t, err, tt.wantErr)
					assert.Equal(t, tt.wantSuccess, success)
				})
			}
		})
	}
}
