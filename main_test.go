package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	t.Parallel()

	executeCMD = func(v string, args []string, outW, err io.Writer) int { return 0 }
	exitCmd = func(i int) {}
	assert.NotPanics(t, func() { main() })
}
