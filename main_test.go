package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	t.Parallel()

	executeCMD = func(s string) int { return 0 }
	exitCmd = func(i int) {}
	assert.NotPanics(t, func() { main() })
}
