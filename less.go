package skiplist

var (
	IntAsc LessFunc = func(a, b interface{}) bool {
		return a.(int) < b.(int)
	}
	IntDesc LessFunc = func(a, b interface{}) bool {
		return a.(int) > b.(int)
	}
	Int64Asc LessFunc = func(a, b interface{}) bool {
		return a.(int64) < b.(int64)
	}
	Int64Desc LessFunc = func(a, b interface{}) bool {
		return a.(int64) > b.(int64)
	}
	Uint64Asc LessFunc = func(a, b interface{}) bool {
		return a.(uint64) < b.(uint64)
	}
	Uint64Desc LessFunc = func(a, b interface{}) bool {
		return a.(uint64) > b.(uint64)
	}
	Int32Asc LessFunc = func(a, b interface{}) bool {
		return a.(int32) < b.(int32)
	}
	Int32Desc LessFunc = func(a, b interface{}) bool {
		return a.(int32) > b.(int32)
	}
	Uint32Asc LessFunc = func(a, b interface{}) bool {
		return a.(uint32) < b.(uint32)
	}
	Uint32Desc LessFunc = func(a, b interface{}) bool {
		return a.(uint32) > b.(uint32)
	}
	StringAsc LessFunc = func(a, b interface{}) bool {
		return a.(string) < b.(string)
	}
	StringDesc LessFunc = func(a, b interface{}) bool {
		return a.(string) > b.(string)
	}
)
