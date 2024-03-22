package main

const n = 4096

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

func main() {
	s := [2]int64{0, 0}
	Add2(s)
}
