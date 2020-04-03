package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	t.Parallel()

	vc := newPathCmd()
	out, _ := prepCmd(t, vc, nil)

	err := vc.Execute()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), pathShort)
	assert.Contains(t, out.String(), pathLong)
}
