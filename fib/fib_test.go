package main

import (
	"testing"
)

// https://stackoverflow.com/questions/36966947/do-go-testing-b-benchmarks-prevent-unwanted-optimizations

var result int

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func BenchmarkFibWrong(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fib(30)
	}
}

func BenchmarkFibBetter(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result = Fib(30)
	}
}

func BenchmarkFibComplete(b *testing.B) {
	var r int
	for n := 0; n < b.N; n++ {
		// always record the result of Fib to prevent
		// the compiler eliminating the function call.
		r = Fib(30)
	}
	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = r
}

// $go test -bench . -benchmem -count=4
// goos: linux
// goarch: amd64
// pkg: github.com/ludwig125/sync-pool/sample2
// BenchmarkFibWrong-8                  213           5415524 ns/op               0 B/op          0 allocs/op
// BenchmarkFibWrong-8                  226           5406527 ns/op               0 B/op          0 allocs/op
// BenchmarkFibWrong-8                  229           5301481 ns/op               0 B/op          0 allocs/op
// BenchmarkFibWrong-8                  225           5410050 ns/op               0 B/op          0 allocs/op
// BenchmarkFibBetter-8                 220           5270236 ns/op               0 B/op          0 allocs/op
// BenchmarkFibBetter-8                 223           5355714 ns/op               0 B/op          0 allocs/op
// BenchmarkFibBetter-8                 219           5210932 ns/op               0 B/op          0 allocs/op
// BenchmarkFibBetter-8                 230           5140387 ns/op               0 B/op          0 allocs/op
// BenchmarkFibComplete-8               216           5381057 ns/op               0 B/op          0 allocs/op
// BenchmarkFibComplete-8               217           5292431 ns/op               0 B/op          0 allocs/op
// BenchmarkFibComplete-8               220           5286499 ns/op               0 B/op          0 allocs/op
// BenchmarkFibComplete-8               235           5390941 ns/op               0 B/op          0 allocs/op
// PASS
// ok      github.com/ludwig125/sync-pool/sample2  20.769s
