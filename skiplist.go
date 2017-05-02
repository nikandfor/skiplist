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
		zero      el
		up        []**el
		autoreuse bool
	}
	el struct {
		val  interface{} /* val */
		h    int
		next [FixedHeight]*el
		more []*el
	}
)

var pool = sync.Pool{New: func() interface{} { return &el{} }}

// New creates skiplist without repeated elements
func New(less LessFunc) *List {
	return &List{
		less:      less,
		zero:      el{h: MaxHeight, more: make([]*el, MaxHeight-FixedHeight)},
		up:        make([]**el, MaxHeight),
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
func (l *List) First() *el {
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
func (l *List) Get(v interface{} /* val */) *el {
	el := l.find(v)
	return el
}
func (l *List) find(v interface{} /* val */) *el {
	cur := &l.zero
loop:
	for {
		// find greatest element that less than v. if any
		for i := cur.height() - 1; i >= 0; i-- {
			next := cur.nexti(i)
			if next == nil {
				continue
			}
			if l.less(next.val, v) {
				cur = next
				continue loop
			}
		}
		next := cur.Next()
		if next != nil && !l.less(v, next.val) {
			cur = next
		}
		// there is no next element less than v
		if cur == &l.zero || l.less(cur.val, v) {
			// no equal elements
			return nil
		}
		return cur
	}
}

// Put puts new value. If it is list with repititions, than it adds new copy after all equals.
// Overwise it rewrites (not replaces) existing.
// It returns positive second argument if there wasn't such element before and vise versa
func (l *List) Put(v interface{} /* val */) (*el, bool) {
	el, ok := l.findPut(v)
	el.val = v
	return el, ok
}
func (l *List) Swap(v interface{} /* val */) (*el, interface{} /* val */, bool) {
	el, ok := l.findPut(v)
	old := el.val
	el.val = v
	return el, old, ok
}
func (l *List) findPut(v interface{} /* val */) (*el, bool) {
	cur := &l.zero
loop:
	for {
		// find greatest element that less than v. if any
		for i := cur.height() - 1; i >= 0; i-- {
			next := cur.nexti(i)
			if next == nil {
				continue
			}
			if !l.less(v, next.val) {
				h := cur.height()
				for i := next.height(); i < h; i++ {
					l.up[i] = cur.nextiaddr(i)
				}

				cur = next

				continue loop
			}
		}
		// there is no next element less than v
		var add bool
		if cur == &l.zero || l.less(cur.val, v) || l.repeat {
			h := cur.height()
			for i := 0; i < h; i++ {
				l.up[i] = cur.nextiaddr(i)
			}

			// add
			l.len++
			add = true
			cur = l.rndEl()
			h = cur.height()
			for i := 0; i < h; i++ {
				cur.setnexti(i, *l.up[i])
				*l.up[i] = cur
			}
		}
		return cur, add
	}
}

func (l *List) GetOrPut(v interface{} /* val */) (*el, bool) {
	el, ok := l.findOrPut(v)
	if ok {
		el.val = v
	}
	return el, ok
}
func (l *List) findOrPut(v interface{} /* val */) (*el, bool) {
	cur := &l.zero
loop:
	for {
		// find greatest element that less than v. if any
		for i := cur.height() - 1; i >= 0; i-- {
			next := cur.nexti(i)
			if next == nil {
				continue
			}
			if l.less(next.val, v) {
				h := cur.height()
				for i := next.height(); i < h; i++ {
					l.up[i] = cur.nextiaddr(i)
				}

				cur = next

				continue loop
			}
		}
		// there is no next element less than v
		next := cur.Next()
		if next != nil && !l.less(v, next.val) {
			return next, false
		}

		var add bool
		if cur == &l.zero || l.less(cur.val, v) {
			h := cur.height()
			for i := 0; i < h; i++ {
				l.up[i] = cur.nextiaddr(i)
			}

			// add
			l.len++
			add = true
			cur = l.rndEl()
			h = cur.height()
			for i := 0; i < h; i++ {
				cur.setnexti(i, *l.up[i])
				*l.up[i] = cur
			}
		}
		return cur, add
	}
}

// Del deletes first occurance equals to v and returns it or nil if it wasn't existed
func (l *List) Del(v interface{} /* val */) *el {
	el := l.findDel(v)
	return el
}
func (l *List) findDel(v interface{} /* val */) *el {
	cur := &l.zero
loop:
	for {
		// find greatest element that less than v. if any
		for i := cur.height() - 1; i >= 0; i-- {
			next := cur.nexti(i)
			if next == nil {
				continue
			}
			if l.less(next.val, v) {
				h := cur.height()
				for i := 0; i < h; i++ {
					l.up[i] = cur.nextiaddr(i)
				}

				h = cur.height()
				for i := next.height(); i < h; i++ {
					l.up[i] = cur.nextiaddr(i)
				}

				cur = next
				continue loop
			}
		}
		prev := cur
		cur = cur.Next()
		// there is no next element less than v
		if cur == nil || l.less(v, cur.val) {
			// didn't have
			return nil
		}

		l.len--

		h := prev.height()
		for i := 0; i < h; i++ {
			l.up[i] = prev.nextiaddr(i)
		}

		h = cur.height()
		for i := 0; i < h; i++ {
			*l.up[i] = cur.nexti(i)
		}

		if l.autoreuse {
			Reuse(cur)
		}

		return cur
	}
}

func (l *List) rndEl() *el {
	h := l.rndHeight()

	e := pool.Get().(*el)
	e.h = h
	if h >= FixedHeight {
		e.more = make([]*el, h-FixedHeight)
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

func (e *el) Value() interface{} /* val */ {
	return e.val
}
func (e *el) Next() *el {
	if e.h == 0 {
		return nil
	}
	return e.next[0]
}
func (e *el) nexti(i int) *el {
	if i < FixedHeight {
		return e.next[i]
	} else {
		return e.more[i-FixedHeight]
	}
}
func (e *el) setnexti(i int, v *el) {
	if i < FixedHeight {
		e.next[i] = v
	} else {
		e.more[i-FixedHeight] = v
	}
}
func (e *el) nextiaddr(i int) **el {
	if i < FixedHeight {
		return &e.next[i]
	} else {
		return &e.more[i-FixedHeight]
	}
}
func (e *el) height() int {
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
func (e *el) String() string {
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
func Reuse(cur *el) {
	cur.more = nil
	pool.Put(cur)
}
