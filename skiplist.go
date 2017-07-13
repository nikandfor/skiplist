package skiplist

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
)

const (
	FixedHeight = 4
	FixedLine   = 5
)

var (
	MaxHeight = 29
)

type (
	LessFunc func(a, b interface{} /* val */) bool
	List     struct {
		less      LessFunc
		repeat    bool
		len       int
		zero      El
		up        []**El
		autoreuse bool
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
	el, _ := l.find(v, true, false, false, false)
	//	el := l.find(v)
	return el
}

// Get returns last occurance of element equal to v (equal defined as !less(e, v) && !less(v, e)) or nil if it doesn't exists.
func (l *List) GetLast(v interface{} /* val */) *El {
	el, _ := l.find(v, false, false, false, false)
	//	el := l.find(v)
	return el
}

// Put puts new value. If it is list with repititions, than it adds new copy after all equals.
// Overwise it rewrites (not replaces) existing.
// Second returned argument is true if there wasn't such element.
func (l *List) Put(v interface{} /* val */) (*El, bool) {
	el, had := l.find(v, false, true, l.repeat, false)
	el.val = v
	return el, !had
}

// Put puts new value. If it is list with repititions, than it adds new copy before all equals.
// Overwise it rewrites (not replaces) existing.
// Second returned argument is true if there wasn't such element.
func (l *List) PutFront(v interface{} /* val */) (*El, bool) {
	el, had := l.find(v, true, true, l.repeat, false)
	el.val = v
	return el, !had
}

// GetOrPut gets first occurance or add new and returns it.
// Second returned argument is true if there wasn't such element.
func (l *List) GetOrPut(v interface{} /* val */) (*El, bool) {
	el, had := l.find(v, true, true, false, false)
	el.val = v
	if !had {
		el.val = v
	}
	return el, !had
}

// Del deletes first occurance equals to v and returns it or nil if it wasn't existed
func (l *List) Del(v interface{} /* val */) *El {
	el, _ := l.find(v, true, false, false, true)
	return el
}

func (l *List) find(v interface{} /* val */, first, add, addm, del bool) (*El, bool) {
	cur := &l.zero
	calcup := add || del

	for {
		var gonext bool
		var next *El
		for i := cur.height() - 1; i >= 0; i-- {
			next = cur.nexti(i)
			if next == nil {
				continue
			}
			if first {
				if l.less(next.val, v) {
					gonext = true
					break
				}
			} else {
				if !l.less(v, next.val) {
					gonext = true
					break
				}
			}
		}
		if !gonext {
			break
		}

		if calcup {
			h := cur.height()
			for i := next.height(); i < h; i++ {
				l.up[i] = cur.nextiaddr(i)
			}
		}
		cur = next
	}

	prev := cur
	var had bool
	if first {
		next := cur.Next()
		if next != nil && !l.less(v, next.val) {
			cur = next
			had = true
		}
	} else if cur != &l.zero && !l.less(cur.val, v) {
		had = true
	}

	if had {
		if !addm && !del {
			return cur, had
		}
	} else {
		if !add {
			return nil, had
		}
	}

	if del {
		h := prev.height()
		for i := 0; i < h; i++ {
			l.up[i] = prev.nextiaddr(i)
		}

		l.len--
		h = cur.height()
		for i := 0; i < h; i++ {
			*l.up[i] = cur.nexti(i)
		}

		if l.autoreuse {
			Reuse(cur)
		}
	} else {
		h := cur.height()
		for i := 0; i < h; i++ {
			l.up[i] = cur.nextiaddr(i)
		}

		// add
		l.len++
		cur = l.rndEl()
		h = cur.height()
		for i := 0; i < h; i++ {
			cur.setnexti(i, *l.up[i])
			*l.up[i] = cur
		}
	}

	return cur, had
}

func (l *List) rndEl() *El {
	h := l.rndHeight()

	e := pool.Get().(*El)
	e.h = h
	if h >= FixedHeight {
		e.more = make([]*El, h-FixedHeight)
	}
	return e
}
func (l *List) rndHeight() int {
	r := rand.Int63()
	h := 1
	for r&1 == 1 && h < len(l.zero.more) {
		h++
		r >>= 1
	}
	return h
}

func (e *El) Value() interface{} /* val */ {
	return e.val
}
func (e *El) Next() *El {
	if e.h == 0 {
		return nil
	}
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
	_, _ = fmt.Fprintf(&buf, "%-10v: (%d)", e.val, e.height())
	for i := 0; i < e.height(); i++ {
		n := e.nexti(i)
		if n == nil {
			_, _ = fmt.Fprintf(&buf, "  nil ")
		} else {
			_, _ = fmt.Fprintf(&buf, "  %-4v", n.val)
		}
	}
	return buf.String()
}

// Put element to buffer for later usage
func Reuse(cur *El) {
	cur.more = nil
	pool.Put(cur)
}
