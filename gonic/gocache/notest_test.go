package gocache

import (
	"github.com/eko/gocache/marshaler"
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
	"github.com/go-redis/redis/v7"

	"github.com/stretchr/testify/assert"
)

var (
	bigcacheClient *bigcache.BigCache
	redisClient *redis.Client

	bigcacheStore *store.BigcacheStore
	redisStore *store.RedisStore
)

func _init() {
	bigcacheClient, _ = bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
	redisClient = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	redisStore = store.NewRedis(redisClient, &store.Options{Expiration: 5*time.Minute})
	bigcacheStore = store.NewBigcache(bigcacheClient, nil) // No otions provided (as second argument)
}

func Test_cacheMgr(t *testing.T) {
	_init()

	mgr := cache.New(bigcacheStore)

	var value = []byte("my-value")
	err := mgr.Set("my-key", value, nil)
	if err != nil {
		panic(err)
	}

	gotValue, _ := mgr.Get("my-key")
	assert.Equal(t, value, gotValue)

	marshal := marshaler.New(mgr)
	// 测试结构体
	type user struct {
		Name string
		Age int
	}

	u := &user{
		Name: "ahah",
		Age: 10,
	}

	err = marshal.Set("user-1", u, nil)
	if err != nil {
		panic(err)
	}
	gotU := new(user)
	_, _ = marshal.Get("user-1", gotU)
	assert.Equal(t, u, gotU)
}

func Test_chain(t *testing.T) {
	_init()

	// Initialize chained cache
	mgr := cache.NewChain(
		cache.New(bigcacheStore),
		cache.New(redisStore),
	)

	mgr.Set("chain-key", []byte("chain-value"), nil)
}