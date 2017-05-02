package skiplist

import "testing"

func TestRepeatedOrder(t *testing.T) {
	MaxHeight = 4
	defer func() {
		MaxHeight = 29
	}()

	type El struct {
		k int
		n int
	}
	l := NewRepeated(func(a, b interface{}) bool {
		return a.(El).k < b.(El).k
	})

	p1, _ := l.Put(El{k: 4, n: 1})
	p2, _ := l.Put(El{k: 4, n: 2})
	p3, _ := l.Put(El{k: 4, n: 3})
	_, _ = l.Put(El{k: 2, n: 3})

	if p1 == nil || p2 == nil || p3 == nil {
		t.Errorf("put 1,2,3: %p %p %p", p1, p2, p3)
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
}
