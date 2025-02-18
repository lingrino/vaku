package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "2.8.3", Version())
}
