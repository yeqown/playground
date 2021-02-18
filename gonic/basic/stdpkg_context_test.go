package basic

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func worker(secs int) (code int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	done := make(chan bool)
	// go func(ctx context.Context) {
	// 	defer fmt.Println("go func defer called")
	// 	fmt.Println("go func statring")

	// 	go func() {
	// 		defer fmt.Println("inner func defer called")
	// 		fmt.Println("inner func statring")
	// 		time.Sleep(time.Duration(secs) * time.Second)
	// 		done <- true
	// 	}()

	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	}
	// }(ctx)

	// type2
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(secs) * time.Second)
			done <- true
		}
	}(ctx)

	// awaiting
	select {
	case <-done:
		// fmt.Println("go func done")
		code = 0
	case <-ctx.Done():
		code = -1
		fmt.Println("timeout")
	}

	return
}

func Test_Worker_context(t *testing.T) {
	type args struct {
		secs int
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		// TODO: Add test cases.
		{name: "case 1", args: args{secs: 3}, wantCode: -1},
		{name: "case 2", args: args{secs: 1}, wantCode: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCode := worker(tt.args.secs); gotCode != tt.wantCode {
				t.Errorf("worker() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}
