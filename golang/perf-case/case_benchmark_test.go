package test

import "testing"

const n = 4097

func Add1(s [2]int64) [2]int64 {
	for i := 0; i < n; i++ {
		s[0]++
		if s[0]%2 == 0 {
			s[1]++
		}
	}
	return s
}

func Add2(s [2]int64) [2]int64 {
	for i := 0; i < n; i++ {
		v := s[0]
		s[0] = v + 1
		if v%2 != 0 {
			s[1]++
		}
	}
	return s
}

func BenchmarkAdd1(b *testing.B) {
	s := [2]int64{0, 0}
	for i := 0; i < b.N; i++ {
		Add1(s)
	}
}

func BenchmarkAdd2(b *testing.B) {
	s := [2]int64{0, 0}
	for i := 0; i < b.N; i++ {
		Add2(s)
	}
}

// go test -bench=. -count=5 -
