package middlewares

import (
	"career-server/logservice"
	"career-server/resource"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 全局限流参数
var (
	windows   = 1000 // 滑动窗口: 1000ms
	limit     = 10   // 默认限流速率
	luaScript = `
		redis.replicate_commands()  -- 关键，切换到命令复制模式

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

func allowRequest(ctx context.Context, key string, windowMs, limit int) (bool, error) {
	res, err := resource.RedisClient.Eval(ctx, luaScript, []string{key}, windowMs, limit).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func GlobalLimiterMiddleware(ctx *gin.Context) {
	key := "global_limit_rate"
	allow, err := allowRequest(ctx, key, windows, limit)
	if err != nil {
		// 如果redis eval exec error，为保证业务正常，不进行限流
		logservice.Notice("redis eval exec error:" + err.Error())
		ctx.Next()
		return
	}
	if !allow {
		logservice.Notice("too many requests")
		ctx.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	ctx.Next()
}
