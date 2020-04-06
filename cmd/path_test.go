package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	t.Parallel()

	vc := newPathCmd()
	stdO, stdE := prepCmd(t, vc, nil)
	assert.Equal(t, "", stdE.String())

	err := vc.Execute()
	assert.NoError(t, err)
	assert.Contains(t, stdO.String(), pathShort)
	assert.Contains(t, stdO.String(), pathLong)
}
