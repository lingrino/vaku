package vaku_test

import (
	"testing"

	"github.com/lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	assert.Equal(t, vaku.Version(), "1.1.1")
}
