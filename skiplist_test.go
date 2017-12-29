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
	defer func(v int) {
		MaxHeight = v
	}(MaxHeight)
	MaxHeight = 5

	l := New(IntLess)

	t.Logf("init:\n%v", l)

	for _, i := range []int{1, 5, 9, 3, 7, 0} {
		cur, add := l.Put(i)
		if !add || cur.Value() == nil || cur.Value().(int) != i {
			t.Errorf("not added: %v %v", i, cur)
		}
	}

	if l.Len() != 6 {
		t.Errorf("Len: %v", l.Len())
	}

	t.Logf("filled:\n%v", l)

	for _, i := range []int{1, 5, 9, 3, 7, 0} {
		cur, add := l.Put(i)
		if add || cur.Value() == nil || cur.Value().(int) != i {
			t.Errorf("added: %v", i)
		}
	}

	if l.Len() != 6 {
		t.Errorf("Len: %v", l.Len())
	}

	t.Logf("filled by the same\n%v", l)

	for _, i := range []int{1, 5, 9, 3, 7, 0} {
		el := l.Get(i)
		if el == nil {
			t.Errorf("Get: %v want %v", el, i)
			continue
		}
		if val, ok := el.Value().(int); !ok || val != i {
			t.Errorf("Get: %v want %v", el, i)
			continue
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
	var prev *El
	for e := l.First(); e != nil; e = e.Next() {
		if e == prev {
			t.Errorf("got after self: %v", e)
			break
		}
		if i >= len(exp) || e.Value() == nil || exp[i] != e.Value().(int) {
			var e int
			if i < len(exp) {
				e = exp[i]
			}
			t.Errorf("at pos %d: %v, want %v", i, e, e)
		}
		i++
		prev = e
	}
	if i < len(exp) {
		t.Errorf("short list: %d", i)
	}

	cur := l.Del(3)
	if cur == nil {
		t.Errorf("%d should be deleted, buf %v", 3, cur)
	}

	t.Logf("del 3\n%v", l)

	if l.Len() != 5 {
		t.Errorf("Len: %v", l.Len())
	}

	cur = l.Del(3)
	if cur != nil {
		t.Errorf("%d already deleted, buf %v", 3, cur)
	}

	t.Logf("del 3 again\n%v", l)

	if l.Len() != 5 {
		t.Errorf("Len: %v", l.Len())
	}

	for _, e := range exp {
		if e == 3 {
			continue
		}
		cur = l.Del(e)
		if cur == nil {
			t.Errorf("should be deleted, buf %v", cur)
		}
	}

	t.Logf("del all\n%v", l)

	if l.Len() != 0 {
		t.Errorf("Len: %v", l.Len())
	}

	for _, e := range exp {
		cur = l.Del(e)
		if cur != nil {
			t.Errorf("already deleted buf %v", cur)
		}
	}
}

func TestPutRepeats(t *testing.T) {
	l := NewRepeated(IntLess)

	l.Put(4)
	l.Put(4)

	l.Put(2)
	l.Put(2)

	l.Put(4)
	l.Put(4)

	exp := []int{2, 2, 4, 4, 4, 4}
	i := 0
	for e := l.First(); e != nil; e = e.Next() {
		if exp[i] != e.Value().(int) {
			t.Errorf("Get: %v want %v", e, exp[i])
		}
		i++
	}
	if i != len(exp) {
		t.Errorf("got %v elements, want %v", i, len(exp))
	}

	t.Logf("filled:\n%v", l)

	l.Del(4)
	l.Del(4)

	l.Del(2)
	l.Del(2)

	l.Del(4)
	l.Del(4)

	t.Logf("del 4:\n%v", l)

	if l.First() != nil {
		t.Errorf("should be 0 elements")
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
		if v == 0 {
			cont := false
			for j := i; j < len(hist) && j < i+3; j++ {
				if hist[j] != 0 {
					cont = true
					break
				}
			}
			if !cont {
				break
			}
		}
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

func TestRandom(t *testing.T) {
	const M = 10000
	l := New(IntLess)

	add := make(map[int]struct{})
	del := make(map[int]struct{})

	for i := 0; i < M; i++ {
		v := rand.Intn(M)
		add[v] = struct{}{}
		l.Put(v)
	}
	if l.Len() != len(add) {
		t.Errorf("Len expected %d, have %d", len(add), l.Len())
	}
	for i := 0; i < M*6/10; i++ {
		v := rand.Intn(M)
		del[v] = struct{}{}
		l.Del(v)
	}

	diff := len(add)
	for v := range add {
		if _, ok := del[v]; ok {
			diff--
			continue
		}

		if el := l.Get(v); el == nil || el.Value() == nil || el.Value().(int) != v {
			t.Errorf("want %d, have %v", v, el)
		}
	}
	if l.Len() != diff {
		t.Errorf("Len expected %d, have %d", diff, l.Len())
	}

	for v := range del {
		if el := l.Get(v); el != nil {
			t.Errorf("want %v, have %v", nil, el)
		}
	}
}

func TestRandomRepeated(t *testing.T) {
	const M = 10000
	l := NewRepeated(IntGreater)
	l2 := NewRepeated(IntGreater)

	add := make(map[int]int)
	del := make(map[int]int)

	for i := 0; i < M; i++ {
		v := rand.Intn(M)
		add[v]++
		l.Put(v)
		l2.GetOrPut(v)
	}
	if l.Len() != M {
		t.Errorf("l.Len expected %d, have %d", M, l.Len())
	}
	if l2.Len() != len(add) {
		t.Errorf("l2.Len expected %d, have %d", len(add), l2.Len())
	}
	for i := 0; i < M*6/10; i++ {
		v := rand.Intn(M)
		del[v]++
		l.Del(v)
		l2.Del(v)
	}

	if M < 50 {
		t.Logf("add: %v", add)
		t.Logf("del: %v", del)
		t.Logf("list:\n%v", l)
		t.Logf("list2:\n%v", l2)
	}

	diff := M
	diff2 := len(add)
	for v, cnt := range add {
		d := del[v]
		if cnt < d {
			d = cnt
		}
		cnt -= d
		diff -= d
		if d > 0 {
			diff2--
		}

		if cnt == 0 {
			if el := l.Get(v); el != nil {
				t.Errorf("want %v, have %v", nil, el)
			}
		} else {
			if el := l.Get(v); el == nil || el.Value() == nil || el.Value().(int) != v {
				t.Errorf("want %d, have %v", v, el)
			}
		}
		if d > 0 {
			if el := l2.Get(v); el != nil {
				t.Errorf("want %v, have %v", nil, el)
			}
		} else {
			if el := l2.Get(v); el == nil || el.Value() == nil || el.Value().(int) != v {
				t.Errorf("want %d, have %v", v, el)
			}
		}
	}
	if l.Len() != diff {
		t.Errorf("Len expected %d, have %d", diff, l.Len())
	}
	if l2.Len() != diff2 {
		t.Errorf("Len expected %d, have %d", diff2, l2.Len())
	}
}

func TestGetOrPut(t *testing.T) {
	l := NewRepeated(IntGreater)

	p1, _ := l.Put(5)
	p2, _ := l.Put(5)
	p3, _ := l.Put(5)

	if p1 == p2 || p1 == p3 {
		t.Errorf("wanted repeats, not the same: %p %p %p", p1, p2, p3)
	}

	if l.Len() != 3 {
		t.Errorf("expected 3 elements, have %v", l.Len())
	}

	g1 := l.Get(5)
	if g1 != p1 {
		t.Errorf("got %p %v, want %p %v", g1, g1, p1, p1)
	}

	gp1, ok := l.GetOrPut(5)
	if ok || gp1 != p1 {
		t.Errorf("got %p %v, want %p %v", gp1, gp1, p1, p1)
	}

	if l.Len() != 3 {
		t.Errorf("expected 3 elements, have %v", l.Len())
	}
}

func BenchmarkAddNewLess(b *testing.B) {
	b.ReportAllocs()

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

func BenchmarkGet(b *testing.B) {
	b.ReportAllocs()

	l := New(IntLess)

	for i := 0; i < b.N; i++ {
		l.Put(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = l.Get(i)
	}
}

func BenchmarkAddNewRepeated(b *testing.B) {
	b.ReportAllocs()

	l := NewRepeated(IntLess)

	for i := 0; i < b.N; i++ {
		l.Put(i)
	}
}
