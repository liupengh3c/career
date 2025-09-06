package middlewares

import (
	"career-server/logservice"
	"career-server/resource"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func Print(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type,token")
	ctx.Header("Access-Control-Max-Age", "2592000")
	ctx.Header("Access-Control-Expose-Headers", "token")
	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusOK)
	}
	start := time.Now()
	ctx.Next()
	msg := fmt.Sprintf(" client_ip[%s],uri[%s] done in %s", ctx.ClientIP(), ctx.Request.RequestURI, time.Since(start))
	logservice.Notice(msg)
}

func getIpLimiter(ip string) *rate.Limiter {
	// 加锁，确保线程安全
	resource.Mux.Lock()
	defer resource.Mux.Unlock()

	// 从全局变量中获取当前IP对应的限流器
	limiter, exists := resource.IpLimiter[ip]
	if !exists {
		// 创建一个新的限流器，设置速率和桶的大小
		limiter = rate.NewLimiter(rate.Limit(resource.IpLimiterCnt), resource.IpLimiterMax)
		// 将新的限流器添加到全局变量中
		resource.IpLimiter[ip] = limiter
	}
	// 返回限流器
	return limiter
}
func IpLimitMiddleware(ctx *gin.Context) {
	// 获取客户端的IP地址
	ip := ctx.ClientIP()
	// 根据IP地址获取访问限制器
	limiter := getIpLimiter(ip)
	// 判断访问限制器是否允许访问
	if !limiter.Allow() {
		// 如果不允许访问，记录日志
		logservice.Notice("too many requests,ip:" + ip)
		// 返回状态码429，表示请求过多
		ctx.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	// 允许访问，继续执行下一个中间件或处理程序
	ctx.Next()
}

func GlobalLimiter(ctx *gin.Context) {
	if !resource.GlobalLimiter.Allow() {
		logservice.Notice("too many requests")
		ctx.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	ctx.Next()
}
