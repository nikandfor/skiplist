package skiplist

import "testing"

func TestCoverLessFuncs(t *testing.T) {
	var l *List

	l = New(IntGreater)
	l.Put(5)
	l.Put(3)

	l = New(IntGE)
	l.Put(5)
	l.Put(3)

	l = New(Int64Less)
	l.Put(int64(5))
	l.Put(int64(3))

	l = New(Int64Greater)
	l.Put(int64(5))
	l.Put(int64(3))

	l = New(Int64LE)
	l.Put(int64(5))
	l.Put(int64(3))

	l = New(Int64GE)
	l.Put(int64(5))
	l.Put(int64(3))

	l = New(Uint64Less)
	l.Put(uint64(5))
	l.Put(uint64(3))

	l = New(Uint64Greater)
	l.Put(uint64(5))
	l.Put(uint64(3))

	l = New(Uint64LE)
	l.Put(uint64(5))
	l.Put(uint64(3))

	l = New(Uint64GE)
	l.Put(uint64(5))
	l.Put(uint64(3))

	l = New(Int32Less)
	l.Put(int32(5))
	l.Put(int32(3))

	l = New(Int32Greater)
	l.Put(int32(5))
	l.Put(int32(3))

	l = New(Int32LE)
	l.Put(int32(5))
	l.Put(int32(3))

	l = New(Int32GE)
	l.Put(int32(5))
	l.Put(int32(3))

	l = New(Uint32Less)
	l.Put(uint32(5))
	l.Put(uint32(3))

	l = New(Uint32Greater)
	l.Put(uint32(5))
	l.Put(uint32(3))

	l = New(Uint32LE)
	l.Put(uint32(5))
	l.Put(uint32(3))

	l = New(Uint32GE)
	l.Put(uint32(5))
	l.Put(uint32(3))

	l = New(StringLess)
	l.Put("5")
	l.Put("3")

	l = New(StringGreater)
	l.Put("5")
	l.Put("3")

	l = New(StringLE)
	l.Put("5")
	l.Put("3")

	l = New(StringGE)
	l.Put("5")
	l.Put("3")
}
