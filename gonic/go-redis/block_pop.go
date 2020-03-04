package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/yeqown/infrastructure/framework/redigo"
	"github.com/yeqown/infrastructure/types"
)

func main() {
	rc, err := redigo.ConnectRedis(&types.RedisConfig{
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "nopass",
	})
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(idx int) {
			defer wg.Done()
			ticker := time.NewTicker(100 * time.Millisecond)
			key := fmt.Sprintf("key-block-%d", idx)
			for {
				select {
				case <-ticker.C:
					log.Printf("%d running before brpop, now=%s\n", idx, time.Now().Format("2006-01-02 15:04:05"))
					res, err := rc.BRPop(5*time.Second, key).Result()
					log.Printf("key=%s, res=%v, err=%v\n", key, res, err)
					log.Printf("%d running after brpop, now=%s\n", idx, time.Now().Format("2006-01-02 15:04:05"))
				default:
					time.Sleep(100 * time.Millisecond)
				}
			}
		}(i)
	}

	wg.Wait()
}
