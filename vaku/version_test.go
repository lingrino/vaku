package vaku_test

import (
	"testing"

	"github.com/Lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, vaku.Version(), "1.0")
}
