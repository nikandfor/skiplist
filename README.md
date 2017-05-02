# go-skiplist
fast skiplist on golang

# Features
* Minimal number of allocs
* Effective + optimized
* It can be used with any types and custom Less function
* Elements can or can not repeats. If elements can repeat, than Get and Del operate on first occurance. Put inserts after all equal elements. (See RepeatedOrder test)
* less than 300 LOC on main file
* There is code generator that replaces `interface{}` to `underlying_type` for even better results
* There are some ready to use Less functions
* tested
* It is invented here

# Benchmarks
```
nik@nik-msi@08:16:08:go-skiplist$ GOMAXPROCS=1 go test . -bench .
BenchmarkAddNewLess 	 2000000	       729 ns/op	      89 B/op	       2 allocs/op
BenchmarkAddDouble  	 2000000	       803 ns/op	      89 B/op	       2 allocs/op
BenchmarkGet        	 5000000	       347 ns/op	       8 B/op	       1 allocs/op
BenchmarkAddNewLE   	 2000000	       824 ns/op	      89 B/op	       2 allocs/op
PASS
ok  	github.com/nikandfor/go-skiplist	16.049s
nik@nik-msi@08:16:31:go-skiplist$ GOMAXPROCS=1 go test . -cover
ok  	github.com/nikandfor/go-skiplist	0.036s	coverage: 96.5% of statements

nik@nik-msi@08:14:43:go-skiplist$ ./make_codegen.sh int
nik@nik-msi@08:16:55:go-skiplist$ go test ./cg/ -bench .
BenchmarkAddNewLess-8   	 3000000	       594 ns/op	      81 B/op	       1 allocs/op
BenchmarkAddDouble-8    	 2000000	       658 ns/op	      81 B/op	       1 allocs/op
BenchmarkGet-8          	 5000000	       261 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddNewLE-8     	 3000000	       591 ns/op	      81 B/op	       1 allocs/op
PASS
ok  	github.com/nikandfor/go-skiplist/cg	13.639s
PASS
```

## Allocs
In `Add` benchmarks one alloc is for a list elements allocation. (but there is sync.Pool in case of you remove elements)
In first group of benchmarks one alloc everywhere is for `int` to `interface{}` convertation.
