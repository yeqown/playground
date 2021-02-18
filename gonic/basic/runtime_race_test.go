package basic_test

import (
	"math/rand"
	"testing"
	"time"
)

// go test -race -run Test_Race -count=1 -v
func Test_Race(t *testing.T) {
	start := time.Now()
	var tm *time.Timer
	tm = time.AfterFunc(randomDuration(), func() {
		t.Logf(time.Now().Sub(start).String())
		tm.Reset(randomDuration())
	})

	time.Sleep(5 * time.Second)
}

func randomDuration() time.Duration {
	return time.Duration(rand.Int63n(1e9))
}
