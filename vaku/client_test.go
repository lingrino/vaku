package vaku

import (
	"testing"

	"github.com/hashicorp/vault/api"

	"github.com/stretchr/testify/assert"
)

const (
	tokenVerifyString = "this token used to verify client equality"
)

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
				dc: &Client{
					workers:      10,
					absolutePath: false,
				},
				workers:      10,
				absolutePath: false,
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
				vc: newDefaultVaultClient(t),
				dc: &Client{
					vc:           newDefaultVaultClient(t),
					workers:      100,
					absolutePath: false,
				},
				workers:      100,
				absolutePath: false,
			},
			wantErr: nil,
		},
		{
			name: "src/dst",
			give: []Option{
				WithVaultSrcClient(newDefaultVaultClient(t)),
				WithVaultDstClient(newDefaultVaultClient(t)),
				WithabsolutePath(true),
			},
			want: &Client{
				vc: newDefaultVaultClient(t),
				dc: &Client{
					vc:           newDefaultVaultClient(t),
					workers:      10,
					absolutePath: true,
				},
				workers:      10,
				absolutePath: true,
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

	if expected.vc != nil {
		assert.Equal(t, expected.vc.Token(), actual.vc.Token())
	} else {
		assert.Nil(t, actual.vc)
	}
	if expected.dc.vc != nil {
		assert.Equal(t, expected.dc.vc.Token(), actual.dc.vc.Token())
	} else {
		assert.Nil(t, actual.vc)
	}

	// zero out clients and assert equal
	expected.vc = nil
	expected.vl = nil
	expected.dc.vc = nil
	expected.dc.vl = nil
	actual.vc = nil
	actual.vl = nil
	actual.dc.vc = nil
	actual.dc.vl = nil

	if expected.dc.dc != expected {
		expected.dc.dc = expected
	}
	if actual.dc.dc != actual {
		actual.dc.dc = actual
	}

	assert.Equal(t, expected.dc, actual.dc)

	expected.dc = nil
	actual.dc = nil
	assert.Equal(t, expected, actual)
}
