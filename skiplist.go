package skiplist

import (
	"bytes"
	"fmt"
	"math/rand"
)

const FixedHeight = 4

var (
	MaxHeight = 29
)

type (
	LessFunc func(a, b interface{}) bool
	List     struct {
		less LessFunc
		zero el
		up   []**el
	}
	el struct {
		val  interface{}
		h    int
		next [FixedHeight]*el
		more []*el
	}
)

var (
	IntLess LessFunc = func(a, b interface{}) bool {
		return a.(int) < b.(int)
	}
)

func New(less LessFunc) *List {
	return &List{
		less: less,
		zero: el{h: MaxHeight, more: make([]*el, MaxHeight-FixedHeight)},
		up:   make([]**el, MaxHeight),
	}
}

func (l *List) First() *el {
	return l.zero.Next()
}

func (l *List) Get(v interface{}) *el {
	el := l.find(v)
	return el
}
func (l *List) find(v interface{}) *el {
	cur := &l.zero
loop:
	for {
		// find greatest element that less than v. if any
		for i := cur.height() - 1; i >= 0; i-- {
			if cur.nexti(i) == nil {
				continue
			}
			if !l.less(v, cur.nexti(i).val) {
				cur = cur.nexti(i)
				continue loop
			}
		}
		//	log.Printf("add after %v", cur)
		// there is no next element less than v
		if cur == &l.zero || l.less(cur.val, v) {
			return nil
		}
		return cur
	}
}

func (l *List) Put(v interface{}) bool {
	el, ok := l.findPut(v)
	el.val = v
	return ok
}
func (l *List) Swap(v interface{}) (interface{}, bool) {
	el, ok := l.findPut(v)
	old := el.val
	el.val = v
	return old, ok
}
func (l *List) findPut(v interface{}) (*el, bool) {
	cur := &l.zero
loop:
	for {
		//	log.Printf("el: %10p %v", cur, cur)
		h := cur.height()
		for i := 0; i < h; i++ {
			l.up[i] = cur.nextiaddr(i)
			//	log.Printf("up[%d]: %p", i, cur.nexti(i))
		}
		// find greatest element that less than v. if any
		for i := cur.height() - 1; i >= 0; i-- {
			if cur.nexti(i) == nil {
				continue
			}
			if !l.less(v, cur.nexti(i).val) {
				cur = cur.nexti(i)
				continue loop
			}
		}
		//	log.Printf("add after %v", cur)
		// there is no next element less than v
		var add bool
		if cur == &l.zero || l.less(cur.val, v) {
			// add
			add = true
			cur = l.rndEl()
			h := cur.height()
			for i := 0; i < h; i++ {
				cur.setnexti(i, *l.up[i])
				*l.up[i] = cur
			}
			//	log.Printf("put %5v %v", v, cur)
		}
		return cur, add
	}
}

func (l *List) Del(v interface{}) bool {
	d := l.findDel(v)
	return d
}
func (l *List) findDel(v interface{}) bool {
	cur := &l.zero
loop:
	for {
		//	log.Printf("el: %10p %v", cur, cur)
		h := cur.height()
		for i := 0; i < h; i++ {
			l.up[i] = cur.nextiaddr(i)
			//	log.Printf("up[%d]: %p", i, cur.nexti(i))
		}
		// find greatest element that less than v. if any
		for i := cur.height() - 1; i >= 0; i-- {
			if cur.nexti(i) == nil {
				continue
			}
			if l.less(cur.nexti(i).val, v) {
				//	if !l.less(v, cur.nexti(i).val) {
				cur = cur.nexti(i)
				continue loop
			}
		}
		cur = cur.Next()
		//	log.Printf("del after %v (%v)", cur, v)
		// there is no next element less than v
		if cur == nil || l.less(v, cur.val) {
			// add
			return false
		}
		h = cur.height()
		for i := 0; i < h; i++ {
			*l.up[i] = cur.nexti(i)
		}
		return true
	}
}

func (l *List) rndEl() *el {
	h := l.rndHeight()

	e := &el{h: h}
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

func (e *el) Value() interface{} {
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
		} else if n.val == nil {
			_, _ = fmt.Fprintf(&buf, "  vnil")
		} else {
			_, _ = fmt.Fprintf(&buf, "  %-4v", n.val)
		}
	}
	return buf.String()
}
