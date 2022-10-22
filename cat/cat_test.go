package cat

import "testing"

func seed(n int) []string {
	s := make([]string, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, "a")
	}
	return s
}

func bench(b *testing.B, n int, f func(...string) string) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f(seed(n)...)
	}
}

func BenchmarkCat3(b *testing.B) {
	bench(b, 3, cat)
}

func BenchmarkBuf3(b *testing.B) {
	bench(b, 3, buf)
}

func BenchmarkCat100(b *testing.B) {
	bench(b, 100, cat)
}

func BenchmarkBuf100(b *testing.B) {
	bench(b, 100, buf)
}

func BenchmarkCat10000(b *testing.B) {
	bench(b, 10000, cat)
}

func BenchmarkBuf10000(b *testing.B) {
	bench(b, 10000, buf)
}
