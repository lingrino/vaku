package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVList(t *testing.T) {
	c := NewClient()
	c.simpleInit()

	l, _ := c.KVList(&KVListInput{
		Path:           "secretv1/test",
		Recurse:        false,
		TrimPathPrefix: true,
	})

	assert.Equal(t, "hi", l)
}
