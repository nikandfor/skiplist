package skiplist

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
)

const (
	FixedHeight = 4
)

var (
	MaxHeight = 30
)

type (
	LessFunc func(a, b interface{} /* val */) bool
	List     struct {
		less      LessFunc
		repeat    bool
		autoreuse bool
		len       int
		zero      El
		up        []**El
	}
	El struct {
		val  interface{} /* val */
		h    int
		next [FixedHeight]*El
		more []*El
	}
)

var pool = sync.Pool{New: func() interface{} { return &El{} }}

// New creates skiplist without repeated elements
func New(less LessFunc) *List {
	return &List{
		less:      less,
		zero:      El{h: MaxHeight, more: make([]*El, MaxHeight-FixedHeight)},
		up:        make([]**El, MaxHeight),
		autoreuse: true,
	}
}

// NewRepeated creates skiplist with possible repeated elements
func NewRepeated(less LessFunc) *List {
	l := New(less)
	l.repeat = true
	return l
}

// First returns first element or nil
func (l *List) First() *El {
	return l.zero.Next()
}

// Len returns length if list
func (l *List) Len() int {
	return l.len
}

// SetAutoReuse enables of disables auto Reuse of deleted elements.
// It is enabled by default.
func (l *List) SetAutoReuse(v bool) {
	l.autoreuse = v
}

// Get returns first occurance of element equal to v (equal defined as !less(e, v) && !less(v, e)) or nil if it doesn't exists.
func (l *List) Get(v interface{} /* val */) *El {
	cur := l.search(v, true, false)

	if cur == nil || l.less(v, cur.val) {
		return nil
	}

	return cur
}

// Get returns last occurance of element equal to v (equal defined as !less(e, v) && !less(v, e)) or nil if it doesn't exists.
func (l *List) GetLast(v interface{} /* val */) *El {
	cur := l.search(v, false, false)

	if cur == &l.zero || l.less(cur.val, v) {
		return nil
	}

	return cur
}

// Put puts new value. If it is list with repititions, than it adds new copy after all equals.
// Overwise it rewrites (not replaces) existing.
// Second returned argument is true if there wasn't such element.
func (l *List) Put(v interface{} /* val */) (*El, bool) {
	cur := l.search(v, false, true)

	if !l.repeat && cur != &l.zero && !l.less(cur.val, v) {
		cur.val = v
		return cur, false
	}

	return l.rndEl(v), true
}

// Put puts new value. If it is list with repititions, than it adds new copy before all equals.
// Overwise it rewrites (not replaces) existing.
// Second returned argument is true if there wasn't such element.
func (l *List) PutBefore(v interface{} /* val */) (*El, bool) {
	cur := l.search(v, true, true)

	if !l.repeat && cur != nil && !l.less(v, cur.val) {
		cur.val = v
		return cur, false
	}

	return l.rndEl(v), true
}

// GetOrPut gets first occurance or add new and returns it.
// Second returned argument is true if there wasn't such element.
func (l *List) GetOrPut(v interface{} /* val */) (*El, bool) {
	cur := l.search(v, true, true)

	if cur != nil && !l.less(v, cur.val) {
		return cur, false
	}

	return l.rndEl(v), true
}

// Del deletes first occurance equals to v and returns it or nil if it wasn't existed
func (l *List) Del(v interface{} /* val */) *El {
	cur := l.search(v, true, true)

	if cur == nil || l.less(v, cur.val) {
		return nil
	}

	l.len--

	h := cur.height()
	for i := h - 1; i >= 0; i-- {
		*l.up[i] = cur.nexti(i)
	}

	if l.autoreuse {
		Reuse(cur)
	}

	return cur
}

func (l *List) DelEl(e *El) *El {
	return l.DelIf(e.Value(), func(b *El) bool { return e == b })
}

func (l *List) DelIf(v interface{} /* val */, f func(*El) bool) *El {
	cur := l.search(v, true, true)

	for cur != nil && !l.less(v, cur.val) && !f(cur) {
		next := cur.Next()
		for i := 0; i < cur.height(); i++ {
			l.up[i] = cur.nextiaddr(i)
		}
		cur = next
	}

	if cur == nil || l.less(v, cur.val) {
		return nil
	}

	l.len--

	h := cur.height()
	for i := h - 1; i >= 0; i-- {
		*l.up[i] = cur.nexti(i)
	}

	if l.autoreuse {
		Reuse(cur)
	}

	return cur
}

func (l *List) search(v interface{} /* val */, first, upd bool) *El {
	cur := &l.zero

	for {
		next := l.jump(cur, v, first, upd)
		if next == nil {
			break
		}

		cur = next
	}

	if first {
		cur = cur.Next()
	}

	return cur
}

func (l *List) jump(cur *El, v interface{} /* val */, first, upd bool) (next *El) {
	for i := cur.height() - 1; i >= 0; i-- {
		n := cur.nexti(i)
		if n == nil {
			continue
		}
		if first {
			if l.less(n.val, v) {
				next = n
				break
			}
		} else {
			if !l.less(v, n.val) {
				next = n
				break
			}
		}
	}
	if upd {
		var nh int
		if next != nil {
			nh = next.height()
		}
		for i := nh; i < cur.height(); i++ {
			l.up[i] = cur.nextiaddr(i)
		}
	}
	return next
}

func (l *List) rndEl(v interface{} /* val */) *El {
	h := l.rndHeight()

	l.len++

	e := pool.Get().(*El)
	e.h = h
	e.val = v
	if h > FixedHeight {
		e.more = make([]*El, h-FixedHeight)
	}

	for i := h - 1; i >= 0; i-- {
		e.setnexti(i, *l.up[i])
		*l.up[i] = e
	}

	return e
}
func (l *List) rndHeight() int {
	r := rand.Int63()
	h := 1
	for r&1 == 1 && h+1 < MaxHeight {
		h++
		r >>= 1
	}
	return h
}

func (e *El) Value() interface{} /* val */ {
	return e.val
}
func (e *El) Next() *El {
	return e.next[0]
}
func (e *El) nexti(i int) *El {
	if i < FixedHeight {
		return e.next[i]
	} else {
		return e.more[i-FixedHeight]
	}
}
func (e *El) setnexti(i int, v *El) {
	if i < FixedHeight {
		e.next[i] = v
	} else {
		e.more[i-FixedHeight] = v
	}
}
func (e *El) nextiaddr(i int) **El {
	if i < FixedHeight {
		return &e.next[i]
	} else {
		return &e.more[i-FixedHeight]
	}
}
func (e *El) height() int {
	return e.h
}

func (l *List) String() string {
	var buf bytes.Buffer
	for z := &l.zero; z != nil; z = z.Next() {
		_, _ = buf.WriteString(z.String())
		_ = buf.WriteByte('\n')
	}
	return buf.String()
}
func (e *El) String() string {
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "%-10v: (%d)", fmt.Sprint(e.val), e.height())
	for i := 0; i < e.height(); i++ {
		n := e.nexti(i)
		if n == nil {
			_, _ = fmt.Fprintf(&buf, "  nil ")
		} else {
			_, _ = fmt.Fprintf(&buf, "  %-4v", fmt.Sprint(n.val))
		}
	}
	return buf.String()
}

// Put element to buffer for later usage
func Reuse(cur *El) {
	cur.more = nil
	pool.Put(cur)
}
