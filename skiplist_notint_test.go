package skiplist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeatedOrder(t *testing.T) {
	defer func(v int) {
		MaxHeight = v
	}(MaxHeight)
	MaxHeight = 4

	type El struct {
		k int
		n int
	}
	l := NewRepeated(func(a, b interface{}) bool {
		return a.(El).k < b.(El).k
	})

	p1, _ := l.Put(El{k: 4, n: 1})
	assert.NotNil(t, p1)
	p2, _ := l.Put(El{k: 4, n: 2})
	assert.NotNil(t, p2)
	p3, _ := l.Put(El{k: 4, n: 3})
	assert.NotNil(t, p3)
	p4, _ := l.Put(El{k: 2, n: 3})
	assert.NotNil(t, p4)

	assert.Equal(t, 4, l.Len())
	if p1 == p2 || p1 == p3 || p2 == p3 {
		t.Fatalf("the same elements: %p %p %p", p1, p2, p3)
	}

	g1 := l.Get(El{k: 4})
	if g1 == nil || g1.Value() == nil || g1.Value().(El).n != 1 {
		t.Fatalf("get: %v", g1)
	}

	d1 := l.Del(El{k: 4})
	if d1 == nil || d1.Value() == nil || d1.Value().(El).n != 1 {
		t.Fatalf("del: %v", d1)
	}

	d2 := l.Del(El{k: 4})
	if d2 == nil || d2.Value() == nil || d2.Value().(El).n != 2 {
		t.Fatalf("del: %v", d2)
	}

	d3 := l.Del(El{k: 4})
	if d3 == nil || d3.Value() == nil || d3.Value().(El).n != 3 {
		t.Fatalf("del: %v", d3)
	}

	f := l.First()
	if f == nil || f.Value() == nil || f.Value().(El).k != 2 {
		t.Fatalf("first: %v", g1)
	}

	assert.Equal(t, 1, l.Len())
}

func TestDelEl(t *testing.T) {
	t.Logf("MaxHeight: %v", MaxHeight)

	type Elt struct {
		k int
		n int
	}
	l := NewRepeated(func(a, b interface{}) bool {
		return a.(Elt).k < b.(Elt).k
	})

	for i := 0; i < 10; i++ {
		l.Put(Elt{k: 1, n: i})
	}

	t.Logf("\n%v", l)

	l.DelCheck(Elt{k: 1}, func(b *El) bool {
		et := b.Value().(Elt)
		return et.n == 2
	})

	i := 0
	for e := l.First(); e != nil; e = e.Next() {
		if i == 2 {
			i++
		}
		assert.Equal(t, i, e.Value().(Elt).n)
		i++
	}

	t.Logf("\n%v", l)

	l.Put(Elt{k: 10})

	l.DelCheck(Elt{k: 1}, func(b *El) bool {
		et := b.Value().(Elt)
		assert.Equal(t, 1, et.k, "stepped over requested element")
		return et.n == 2
	})
}
