package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"

	vault "github.com/hashicorp/vault/api"
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
				WithAbsolutePath(true),
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
			name:    "bad workers",
			give:    []Option{WithWorkers(0)},
			want:    nil,
			wantErr: []error{ErrNumWorkers},
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

func TestSwapPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		giveSrc       string
		giveDst       string
		giveAbsData   map[string]map[string]interface{}
		giveNoAbsData map[string]map[string]interface{}
		wantAbs       map[string]map[string]interface{}
		wantNoAbs     map[string]map[string]interface{}
	}{
		{
			giveSrc: "0/1/2",
			giveDst: "00/01/02",
			giveAbsData: map[string]map[string]interface{}{
				"0/1/2/3": nil,
				"0/1/2/4": nil,
			},
			giveNoAbsData: map[string]map[string]interface{}{
				"0/1/2/3": nil,
				"0/1/2/4": nil,
			},
			wantAbs: map[string]map[string]interface{}{
				"00/01/02/3": nil,
				"00/01/02/4": nil,
			},
			wantNoAbs: map[string]map[string]interface{}{
				"00/01/02/0/1/2/3": nil,
				"00/01/02/0/1/2/4": nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.giveSrc, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(WithAbsolutePath(true))
			assert.NoError(t, err)

			client.swapPaths(tt.giveAbsData, tt.giveSrc, tt.giveDst)
			assert.Equal(t, tt.wantAbs, tt.giveAbsData)
		})
		t.Run(tt.giveSrc, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient()
			assert.NoError(t, err)

			client.swapPaths(tt.giveNoAbsData, tt.giveSrc, tt.giveDst)
			assert.Equal(t, tt.wantNoAbs, tt.giveNoAbsData)
		})
	}
}

func TestOutputPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		giveRoot  string
		givePath  string
		wantAbs   string
		wantNoAbs string
	}{
		{
			giveRoot:  "0/1/2",
			givePath:  "3",
			wantAbs:   "0/1/2/3",
			wantNoAbs: "3",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.givePath, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(WithAbsolutePath(true))
			assert.NoError(t, err)

			res := client.outputPath(tt.givePath, tt.giveRoot)
			assert.Equal(t, tt.wantAbs, res)
		})
		t.Run(tt.givePath, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient()
			assert.NoError(t, err)

			res := client.outputPath(tt.givePath, tt.giveRoot)
			assert.Equal(t, tt.wantNoAbs, res)
		})
	}
}

func TestOutputPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		giveRoot       string
		giveAbsPaths   []string
		giveNoAbsPaths []string

		wantAbs   []string
		wantNoAbs []string
	}{
		{
			giveRoot:       "0/1/2",
			giveAbsPaths:   []string{"3", "4"},
			giveNoAbsPaths: []string{"3", "4"},
			wantAbs:        []string{"0/1/2/3", "0/1/2/4"},
			wantNoAbs:      []string{"3", "4"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.giveRoot, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(WithAbsolutePath(true))
			assert.NoError(t, err)

			client.outputPaths(tt.giveAbsPaths, tt.giveRoot)
			assert.Equal(t, tt.wantAbs, tt.giveAbsPaths)
		})
		t.Run(tt.giveRoot, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient()
			assert.NoError(t, err)

			client.outputPaths(tt.giveNoAbsPaths, tt.giveRoot)
			assert.Equal(t, tt.wantNoAbs, tt.giveNoAbsPaths)
		})
	}
}

// withError returns the passed in error for Option error injection.
func withError(e error) Option {
	return withErrorOpt{e}
}

type withErrorOpt struct {
	err error
}

func (o withErrorOpt) apply(c *Client) error {
	return o.err
}

// newDefaultVaultClient creates a default vault client and fails on error.
func newDefaultVaultClient(t *testing.T) *vault.Client {
	t.Helper()

	client, err := vault.NewClient(vault.DefaultConfig())
	assert.NoError(t, err)
	client.SetToken(tokenVerifyString)

	return client
}

// assertClientsEqual compares two Clients.
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
