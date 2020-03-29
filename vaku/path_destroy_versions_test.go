package vaku_test

import (
	"fmt"
	"testing"

	"github.com/lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

type TestPathDestroyVersionsData struct {
	input              *vaku.PathInput
	destroyVersionsErr bool
	readErr            bool
}

func TestPathDestroyVersions(t *testing.T) {
	var err error

	c := clientInitForTests(t)

	defer func() {
		err = seed(t, c)
		if err != nil {
			t.Error(fmt.Errorf("failed to reseed: %w", err))
		}
	}()

	tests := map[int]TestPathDestroyVersionsData{
		1: {
			input:              vaku.NewPathInput("secretv1/test/foo"),
			destroyVersionsErr: true,
		},
		2: {
			input:              vaku.NewPathInput("secretv2/test/foo"),
			destroyVersionsErr: false,
			readErr:            false,
		},
		3: {
			input:              vaku.NewPathInput("secretv1/doesnotexist"),
			destroyVersionsErr: true,
		},
		4: {
			input:              vaku.NewPathInput("secretv2/doesnotexist"),
			destroyVersionsErr: false,
			readErr:            true,
		},
		5: {
			input:              vaku.NewPathInput("secretdoesnotexist/test/foo"),
			destroyVersionsErr: true,
		},
	}

	for _, d := range tests {
		destroyVersionsErr := c.PathDestroyVersions(d.input, []int{1})
		_, readErr := c.PathRead(d.input)
		if d.destroyVersionsErr {
			assert.Error(t, destroyVersionsErr)
		} else {
			assert.NoError(t, destroyVersionsErr)
			if d.readErr {
				assert.Error(t, readErr)
			} else {
				assert.NoError(t, readErr)
			}
		}
	}
}
