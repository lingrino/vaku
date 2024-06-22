package vaku

import (
	"github.com/chmduquesne/rollinghash/buzhash32"
)

const defaultWindowSize = 10

// haystacksearch package content
var (
	// DefaultHashMaker is the default hash maker used in haystacksearch. It is backed by buzhash32.
	DefaultHashMaker = NewRollingHashMaker()
)

// RollingHashMaker defines an interface that produces RollingHashers.
type RollingHashMaker interface {
	Make() RollingHasher
}

type buzhashMaker struct{}

func (b *buzhashMaker) Make() RollingHasher {
	return newRollingHasher()
}

// NewRollingHashMaker creates a new RollingHashMaker backed by the buzhash32 rolling hash algorithm.
func NewRollingHashMaker() RollingHashMaker {
	return &buzhashMaker{}
}

// RollingHasher defines an interface for performing rolling hash writes and summation.
type RollingHasher interface {
	// Write writes data to the buffer and returns the number of bytes written.
	Write(data []byte) (int, error)

	// Sum32 returns the current sum of the buffer.
	Sum32() uint32

	// Roll adds a byte to the current buffer.
	Roll(c byte)
}

type rollingBuzhasher struct {
	hasher *buzhash32.Buzhash32
}

// newRollingHasher creates a new RollingHasher backed by the buzhash32 rolling hash algorithm.
func newRollingHasher() RollingHasher {
	return &rollingBuzhasher{hasher: buzhash32.New()}
}

// Write writes data to the buffer and returns the number of bytes written.
func (r *rollingBuzhasher) Write(data []byte) (int, error) {
	return r.hasher.Write(data)
}

// Sum32 returns the current sum of the buffer.
func (r *rollingBuzhasher) Sum32() uint32 {
	return r.hasher.Sum32()
}

// Roll adds a byte to the current hash.
func (r *rollingBuzhasher) Roll(c byte) {
	r.hasher.Roll(c)
}
