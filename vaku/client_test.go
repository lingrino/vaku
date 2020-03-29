package vaku

import (
	"testing"

	"github.com/hashicorp/vault/api"

	"github.com/stretchr/testify/assert"
)

const (
	tokenVerifyString = "this token used to verify client equality"
)

// withError returns the passed in error for Option error injection
func withError(e error) Option {
	return withErrorOpt{e}
}

type withErrorOpt struct {
	err error
}

func (o withErrorOpt) apply(c *Client) error {
	return o.err
}

// newDefaultVaultClient creates a default vault client and fails on error
func newDefaultVaultClient(t *testing.T) *api.Client {
	t.Helper()

	client, err := api.NewClient(api.DefaultConfig())
	assert.NoError(t, err)
	client.SetToken(tokenVerifyString)

	return client
}

// assertClientsEqual compares two Clients
func assertClientsEqual(t *testing.T, expected *Client, actual *Client) {
	if expected == nil {
		assert.Nil(t, actual)
		return
	}

	if expected.src != nil {
		assert.Equal(t, expected.src.Token(), actual.src.Token())
	} else {
		assert.Nil(t, actual.src)
	}
	if expected.dst != nil {
		assert.Equal(t, expected.dst.Token(), actual.dst.Token())
	} else {
		assert.Nil(t, actual.dst)
	}

	// zero out clients and assert equal
	expected.src = nil
	expected.srcL = nil
	expected.dst = nil
	expected.dstL = nil
	actual.src = nil
	actual.srcL = nil
	actual.dst = nil
	actual.dstL = nil
	assert.Equal(t, expected, actual)
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    []Option
		want    *Client
		wantErr []error
	}{
		{
			name: "nil",
			give: []Option{},
			want: &Client{
				workers: 10,
			},
			wantErr: nil,
		},
		{
			name: "vault client",
			give: []Option{
				WithVaultClient(newDefaultVaultClient(t)),
				WithWorkers(100),
			},
			want: &Client{
				src:     newDefaultVaultClient(t),
				dst:     newDefaultVaultClient(t),
				workers: 100,
			},
			wantErr: nil,
		},
		{
			name: "src/dst",
			give: []Option{
				WithVaultSrcClient(newDefaultVaultClient(t)),
				WithVaultDstClient(newDefaultVaultClient(t)),
			},
			want: &Client{
				src:     newDefaultVaultClient(t),
				dst:     newDefaultVaultClient(t),
				workers: 10,
			},
			wantErr: nil,
		},
		{
			name: "error",
			give: []Option{
				withError(errInject),
			},
			want:    nil,
			wantErr: []error{errInject},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(tt.give...)

			compareErrors(t, err, tt.wantErr)
			assertClientsEqual(t, tt.want, client)
		})
	}
}
