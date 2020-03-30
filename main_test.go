package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	executeCMD = func(s string) {}
	assert.NotPanics(t, func() { main() })
}
