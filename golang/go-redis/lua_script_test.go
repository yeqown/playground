package main

import (
	"testing"
	"time"

	"github.com/yeqown/infrastructure/framework/redigo"
	"github.com/yeqown/infrastructure/types"
)

// lua 脚本提供原子操作和并发安全
var script = `
local ok = redis.call('setnx', KEYS[1], ARGV[1]) 
if (not ok) then
	redis.call('del', KEYS[1]) 
	return false 
end
local cnt = redis.call('get', KEYS[2])
if (cnt <= ARGV[2]) then 
	redis.call('del', KEYS[1]) 
	return false 
end
redis.call('decr', KEYS[2]) 
redis.call('del', KEYS[1]) 
return true
`

// 展示如何使用lua脚本来实现原子操作
func Test_LuaScript(t *testing.T) {
	client, err := redigo.ConnectRedis(&types.RedisConfig{
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "nopass",
	})
	if err != nil {
		panic(err)
	}

	keys := []string{"lua_lock0", "lua_cnt"}
	args := []interface{}{"locked", 0}

	// 准备数据
	if err := client.
		Set(keys[1], 10, 10*time.Second).Err(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// 上传脚本
	// return slice of return values
	cmd := client.Eval(script, keys, args...)
	if cmd.Err() != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(cmd.Val()) // result values

	time.Sleep(100 * time.Millisecond)

	// `script load $SCRIPT`
	scriptSha := "02da95915a026780ab9e8e775cbed99fcd00755c"
	t.Log("sha of script is: ", scriptSha)

	cmd2 := client.EvalSha(scriptSha, keys, args...)
	if cmd.Err() != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(cmd2.Val()) // result values
}
