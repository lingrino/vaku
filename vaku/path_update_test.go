package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathUpate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        string
		giveData    map[string]interface{}
		giveLogical logical
		giveOptions []Option
		wantData    map[string]interface{}
		wantErr     []error
	}{
		{
			name: "new path",
			give: "update/bar",
			giveData: map[string]interface{}{
				"Eg5ljS7t": "6F1B5nBg",
				"quqr32S5": "81iY4HAN",
				"r6R0JUzX": "rs1mCRB5",
			},
			wantData: map[string]interface{}{
				"Eg5ljS7t": "6F1B5nBg",
				"quqr32S5": "81iY4HAN",
				"r6R0JUzX": "rs1mCRB5",
			},
			wantErr: nil,
		},
		{
			name: "existing path",
			give: "test/foo",
			giveData: map[string]interface{}{
				"foo": "bar",
			},
			wantData: map[string]interface{}{
				"foo":   "bar",
				"value": "bar",
			},
			wantErr: nil,
		},
		{
			name: "partial overwrite",
			give: "test/value",
			giveData: map[string]interface{}{
				"fizz": "bar",
			},
			wantData: map[string]interface{}{
				"fizz": "bar",
				"foo":  "bar",
			},
			wantErr: nil,
		},
		{
			name:     "nil data new path",
			give:     "update/nildata",
			giveData: nil,
			wantErr:  []error{ErrPathUpdate, ErrNilData},
		},
		{
			name:     "nil data existing path",
			give:     "test/foo",
			giveData: nil,
			wantData: map[string]interface{}{
				"value": "bar",
			},
			wantErr: []error{ErrPathUpdate, ErrNilData},
		},
		{
			name: "no mount",
			give: noMountPrefix,
			giveData: map[string]interface{}{
				"foo":   "bar",
				"value": "bar",
			},
			wantErr: []error{ErrPathUpdate, ErrPathWrite, ErrVaultWrite},
		},
		{
			name: "inject read",
			give: "test/foo",
			giveData: map[string]interface{}{
				"foo":   "bar",
				"value": "bar",
			},
			giveLogical: &errLogical{
				err: errInject,
				op:  "Read",
			},
			wantData: map[string]interface{}{
				"value": "bar",
			},
			wantErr: []error{ErrPathUpdate, ErrPathRead, ErrVaultRead},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := testClient(t, tt.giveOptions...)
			readbackClient := cloneCLient(t, client)
			updateLogical(t, client, tt.giveLogical, tt.giveLogical)

			for _, ver := range kvMountVersions {
				path := addMountToPath(t, tt.give, ver)

				err := client.PathUpdate(path, tt.giveData)
				compareErrors(t, err, tt.wantErr)

				readBack, err := readbackClient.PathRead(path)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantData, readBack)
			}
		})
	}
}
