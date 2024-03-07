package main

const n = 4096

func Add1(s [2]int64) [2]int64 {
	for i := 0; i < n; i++ {
		s[0]++
		if s[0]%2 == 0 {
			s[1]++
		}
	}
	return s
}

func main() {
	s := [2]int64{0, 0}
	Add1(s)
}

// go tool compile -S main.go -N -l > main.s
