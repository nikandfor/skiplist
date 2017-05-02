package skiplist

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
)

const FixedHeight = 4

var (
	MaxHeight = 29
)

type (
	LessFunc func(a, b interface{} /* val */) bool
	List     struct {
		less LessFunc
		eq   LessFunc
		zero el
		up   []**el
	}
	el struct {
		val  interface{} /* val */
		h    int
		next [FixedHeight]*el
		more []*el
	}
)

var pool = sync.Pool{New: func() interface{} { return &el{} }}

func New(less LessFunc) *List {
	return &List{
		less: less,
		zero: el{h: MaxHeight, more: make([]*el, MaxHeight-FixedHeight)},
		up:   make([]**el, MaxHeight),
	}
}
func NewLE(less LessFunc) *List {
	l := New(less)
	l.eq = func(a, b interface{} /* val */) bool {
		return less(b, a)
	}
	return l
}

func (l *List) First() *el {
	return l.zero.Next()
}

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
			if !l.less(v, next.val) {
				cur = next
				continue loop
			}
		}
		// there is no next element less than v
		if cur == &l.zero || l.less(cur.val, v) {
			// no equal elements
			return nil
		}
		return cur
	}
}

func (l *List) Put(v interface{} /* val */) bool {
	el, ok := l.findPut(v)
	el.val = v
	return ok
}
func (l *List) Swap(v interface{} /* val */) (interface{} /* val */, bool) {
	el, ok := l.findPut(v)
	old := el.val
	el.val = v
	return old, ok
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
		if cur == &l.zero || l.less(cur.val, v) {

			h := cur.height()
			for i := 0; i < h; i++ {
				l.up[i] = cur.nextiaddr(i)
			}

			// add
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

func (l *List) Del(v interface{} /* val */) bool {
	d := l.findDel(v)
	return d
}
func (l *List) findDel(v interface{} /* val */) bool {
	cur := &l.zero
loop:
	for {
		// find greatest element that less than v. if any
		for i := cur.height() - 1; i >= 0; i-- {
			next := cur.nexti(i)
			if next == nil {
				continue
			}
			if l.less(next.val, v) && (l.eq == nil || !l.eq(next.val, v)) {
				h := cur.height()
				for i := 0; i < h; i++ {
					l.up[i] = cur.nextiaddr(i)
					//	log.Printf("up[%d]: %p", i, cur.nexti(i))
				}

				//	if !l.less(v, cur.nexti(i).val) {

				h = cur.height()
				for i := next.height(); i < h; i++ {
					l.up[i] = cur.nextiaddr(i)
					//	log.Printf("up[%d]: %p", i, cur.nexti(i))
				}

				cur = next
				continue loop
			}
		}
		cur = cur.Next()
		// there is no next element less than v
		if cur == nil || l.less(v, cur.val) && (l.eq == nil || !l.eq(v, cur.val)) {
			// didn't have
			return false
		}
		h = cur.height()
		for i := 0; i < h; i++ {
			*l.up[i] = cur.nexti(i)
		}
		cur.more = nil
		pool.Put(cur)
		return true
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
