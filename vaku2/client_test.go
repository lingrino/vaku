package vaku2

import (
	"errors"
	"testing"

	"github.com/hashicorp/vault/api"

	"github.com/stretchr/testify/assert"
)

const (
	tokenVerifyString = "this token used to verify client equality"
)

var (
	errTestNewClient = errors.New("test error when creating a new client")
)

type withErrorOpt struct {
	err error
}

func (o withErrorOpt) apply(c *Client) error {
	return o.err
}

// withError returns the passed in error
func withError(e error) Option {
	return withErrorOpt{e}
}

// newDefaultVaultClient creates a default vault client and fails on error
func newDefaultVaultClient(t *testing.T) *api.Client {
	t.Helper()

	client, err := api.NewClient(api.DefaultConfig())
	assert.NoError(t, err)
	client.SetToken(tokenVerifyString)

	return client
}

// assertClientsEqual compares two Client
func assertClientsEqual(t *testing.T, expected *Client, actual *Client) {
	if expected == nil {
		assert.Nil(t, actual)
		return
	}

	if expected.source != nil {
		assert.Equal(t, expected.source.Token(), actual.source.Token())
	} else {
		assert.Nil(t, actual.source)
	}
	if expected.dest != nil {
		assert.Equal(t, expected.dest.Token(), actual.dest.Token())
	} else {
		assert.Nil(t, actual.dest)
	}
}

// TestNewClient tests NewClient.
func TestNewClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    []Option
		want    *Client
		wantErr error
	}{
		{
			name:    "nil",
			give:    []Option{},
			want:    &Client{},
			wantErr: nil,
		},
		{
			name: "vault client",
			give: []Option{
				WithVaultClient(newDefaultVaultClient(t)),
			},
			want: &Client{
				source: newDefaultVaultClient(t),
			},
			wantErr: nil,
		},
		{
			name: "source/dest",
			give: []Option{
				WithVaultSourceClient(newDefaultVaultClient(t)),
				WithVaultDestClient(newDefaultVaultClient(t)),
			},
			want: &Client{
				source: newDefaultVaultClient(t),
				dest:   newDefaultVaultClient(t),
			},
			wantErr: nil,
		},
		{
			name: "error",
			give: []Option{
				withError(errTestNewClient),
			},
			want:    nil,
			wantErr: errTestNewClient,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(tt.give...)

			// assert the outputs
			assert.Equal(t, tt.wantErr, err)
			assertClientsEqual(t, tt.want, client)
		})
	}
}
