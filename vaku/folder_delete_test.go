package vaku_test

import (
	"fmt"
	"testing"

	"github.com/lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

type TestFolderDeleteData struct {
	input     *vaku.PathInput
	outputErr bool
}

func TestFolderDelete(t *testing.T) {
	var err error

	c := clientInitForTests(t)

	defer func() {
		err = seed(t, c)
		if err != nil {
			t.Error(fmt.Errorf("failed to reseed: %w", err))
		}
	}()

	tests := map[int]TestFolderDeleteData{
		1: {
			input:     vaku.NewPathInput("secretv1/test"),
			outputErr: false,
		},
		2: {
			input:     vaku.NewPathInput("secretv2/test"),
			outputErr: false,
		},
		3: {
			input:     vaku.NewPathInput("secretdoesnotexist/test"),
			outputErr: true,
		},
	}

	for _, d := range tests {
		e := c.FolderDelete(d.input)
		r, re := c.FolderRead(d.input)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			if re == nil {
				for _, data := range r {
					assert.Equal(t, "SECRET_HAS_BEEN_DELETED", data["VAKU_STATUS"])
				}
			} else {
				assert.Error(t, re)
			}
			assert.NoError(t, e)
		}
	}
}
