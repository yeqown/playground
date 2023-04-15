package singleflight_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/singleflight"
)

var (
	group singleflight.Group
	// cache    *sync.Map
	// cacheKey string
)

func init() {
	group = singleflight.Group{}
	// cache = &sync.Map{}
	// cacheKey = "cacheKey"
}

func getFromCache() (interface{}, error) {
	// 从缓存读取
	return nil, errors.New("not hit")
}

// 模拟缓存击穿的场景
func getDB() (interface{}, error) {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("call getDB")
	// 数据库结果返回
	return "db", nil
}

func Test_Nomarl(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			var (
				// v   interface{}
				err error
			)

			if _, err = getFromCache(); err != nil {
				_, _ = getDB()
			}
			// t.Logf("Test_Nomarl got result=%s", v.(string))
		}()
	}

	wg.Wait()
}

func Test_GroupDo(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(10)

	var key = "key"

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			var (
				v      interface{}
				err    error
				shared bool
			)

			if v, err = getFromCache(); err != nil {
				// true: 缓存未命中
				v, _, shared = group.Do(key, getDB)
			}

			t.Logf("Test_GroupDo got result=%s, shared=%v", v.(string), shared)
		}()
	}
	wg.Wait()
}

func Test_GroupDoChan(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(10)

	var keyChan = "keyChan"

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			var (
				v      interface{}
				err    error
				shared bool
			)

			if v, err = getFromCache(); err != nil {
				// true: 缓存未命中
				res := <-group.DoChan(keyChan, getDB)
				v, _, shared = res.Val, res.Err, res.Shared
			}

			t.Logf("Test_GroupDo got result=%s, shared=%v", v.(string), shared)
		}()
	}
	wg.Wait()
}
