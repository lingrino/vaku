package vaku2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	assert.Equal(t, Version(), "2.0.0")
}
