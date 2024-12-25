package binary_test

import (
	"testing"

	"github.com/daanv2/go-cache/pkg/binary"
	"github.com/stretchr/testify/assert"
)

func Test_Uint64(t *testing.T) {
	assert.EqualValues(t, binary.Uint64([]byte{}), 0)
	assert.EqualValues(t, binary.Uint64([]byte{1}), 0x1)
	assert.EqualValues(t, binary.Uint64([]byte{1, 2}), 0x201)
	assert.EqualValues(t, binary.Uint64([]byte{1, 2, 3}), 0x30201)
	assert.EqualValues(t, binary.Uint64([]byte{1, 2, 3, 4}), 0x4030201)
	assert.EqualValues(t, binary.Uint64([]byte{1, 2, 3, 4, 5}), 0x504030201)
	assert.EqualValues(t, binary.Uint64([]byte{1, 2, 3, 4, 5, 6}), 0x60504030201)
	assert.EqualValues(t, binary.Uint64([]byte{1, 2, 3, 4, 5, 6, 7}), 0x7060504030201)
	assert.EqualValues(t, binary.Uint64([]byte{1, 2, 3, 4, 5, 6, 7, 8}), 0x807060504030201)
	assert.EqualValues(t, binary.Uint64([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}), 0x807060504030201)
}
