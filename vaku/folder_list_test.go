package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        string
		giveLogical logical
		giveOptions []Option
		want        []string
		wantErr     []error
	}{
		{
			name:    "test/foo",
			give:    "test/foo",
			want:    nil,
			wantErr: nil,
		},
		{
			name: "test/inner",
			give: "test/inner",
			want: []string{
				"WKNC3muM",
				"A2xlzTfE",
				"again/inner/UCrt6sZT",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ln, client := testClient(t, tt.giveOptions...)
			defer ln.Close()
			updateLogical(t, client, tt.giveLogical, tt.giveLogical)

			for _, ver := range kvMountVersions {
				path := addMountToPath(t, tt.give, ver)

				list, err := client.FolderList(path)
				compareErrors(t, err, tt.wantErr)

				assert.Equal(t, tt.want, list)
			}
		})
	}
}
