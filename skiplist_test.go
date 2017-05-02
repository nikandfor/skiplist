package skiplist

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestPutGet(t *testing.T) {
	l := New(IntLess)

	t.Logf("\n%v", l)

	for _, i := range []int{1, 5, 9, 3, 7, 0} {
		add := l.Put(i)
		if !add {
			t.Errorf("not added: %v", i)
		}
	}

	t.Logf("\n%v", l)

	for _, i := range []int{1, 5, 9, 3, 7, 0} {
		add := l.Put(i)
		if add {
			t.Errorf("added: %v", i)
		}
	}

	t.Logf("\n%v", l)

	for _, i := range []int{1, 5, 9, 3, 7, 0} {
		el := l.Get(i)
		if val, ok := el.Value().(int); !ok || val != i {
			t.Errorf("Get: %v want %v", el, i)
		}
	}

	for _, i := range []int{-1, 2, 4, 6, 8, 10, 300} {
		el := l.Get(i)
		if el != nil {
			t.Errorf("Get: %v want %v", el, nil)
		}
	}

	exp := []int{0, 1, 3, 5, 7, 9}
	i := 0
	var prev *el
	for e := l.First(); e != nil; e = e.Next() {
		if e == prev {
			t.Errorf("got after self: %v", e)
			break
		}
		if i >= len(exp) || e.Value() == nil || exp[i] != e.Value().(int) {
			t.Errorf("at pos %d: %v, want %v", i, e, exp[i])
		}
		i++
		prev = e
	}
	if i < len(exp) {
		t.Errorf("short list: %d", i)
	}
}

func TestHeight(t *testing.T) {
	l := New(nil)
	const D = 1.8
	const Min = 50

	hist := make([]int, 40)
	for i := 0; i < 1000000; i++ {
		h := l.rndHeight()
		hist[0]++
		hist[h]++
	}

	for i, v := range hist {
		p := 0.5
		if i != 0 {
			p = float64(v) / float64(hist[i-1])
		}
		if v > Min && (p > 0.5*D || p < 0.5/D) {
			t.Errorf("i %2d: %7v (%.2f)  <- out of (%.3v %.3v)", i, v, p, 0.5/D, 0.5*D)
		} else {
			t.Logf("i %2d: %7v (%.2f)", i, v, p)
		}
	}
}

func BenchmarkAddNew(b *testing.B) {
	b.ReportAllocs()

	//	l := New(func(a, b interface{}) bool { return rand.Int()%2 == 0 })
	l := New(IntLess)

	for i := 0; i < b.N; i++ {
		l.Put(i)
	}
}

func BenchmarkAddDouble(b *testing.B) {
	b.ReportAllocs()

	l := New(IntLess)

	for i := 0; i < b.N; i++ {
		l.Put(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Put(b.N + i)
	}
}
