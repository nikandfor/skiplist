package skiplist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoverLessFuncs(t *testing.T) {
	assert.True(t, IntLess(1, 2))
	assert.True(t, IntGreater(2, 1))

	assert.True(t, Int64Less(int64(1), int64(2)))
	assert.True(t, Int64Greater(int64(2), int64(1)))

	assert.True(t, Uint64Less(uint64(1), uint64(2)))
	assert.True(t, Uint64Greater(uint64(2), uint64(1)))

	assert.True(t, Int32Less(int32(1), int32(2)))
	assert.True(t, Int32Greater(int32(2), int32(1)))

	assert.True(t, Uint32Less(uint32(1), uint32(2)))
	assert.True(t, Uint32Greater(uint32(2), uint32(1)))

	assert.True(t, StringLess("1", "2"))
	assert.True(t, StringGreater("2", "1"))
}
