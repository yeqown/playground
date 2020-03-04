package golang

import (
	"context"
	"fmt"
	"time"
)

func Worker(secs int) (code int) {
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
