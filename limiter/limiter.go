package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// 全局限流参数
var (
	redisClient *redis.Client
	windows     = 1000 // 滑动窗口: 1000ms
	limit       = 10   // 默认限流速率
	luaScript   = `
		redis.replicate_commands()  -- 关键

		local key = KEYS[1]
		local window = tonumber(ARGV[1])
		local limit = tonumber(ARGV[2])
		
		local timeData = redis.call("TIME")
		local now = timeData[1] * 1000 + math.floor(timeData[2] / 1000)
		
		-- 删除过期请求
		redis.call("ZREMRANGEBYSCORE", key, 0, now - window)
		
		-- 当前请求数
		local count = redis.call("ZCARD", key)
		
		if count < limit then
			redis.call("ZADD", key, now, now)
			redis.call("PEXPIRE", key, window)
			return 1
		else
			return 0
		end
	`
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
func allowRequest(ctx context.Context, key string, windowMs, limit int) (bool, error) {
	res, err := redisClient.Eval(ctx, luaScript, []string{key}, windowMs, limit).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func GlobalLimit() bool {
	key := "global_limit_rate"
	allow, err := allowRequest(context.Background(), key, windows, limit)
	if err != nil {
		return false
	}
	return allow
}

func main() {
	allow := GlobalLimit()
	fmt.Println(allow)
}
