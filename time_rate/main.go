package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// 创建限流器：每秒 2 个请求，突发容量 5
	limiter := rate.NewLimiter(1, 5)

	// 模拟发送 10 个 HTTP 请求
	for i := 0; i < 50; i++ {
		// 使用 Wait 阻塞等待令牌
		if err := limiter.Wait(context.Background()); err != nil {
			fmt.Println("Rate limit error:", err)
			continue
		}

		// 发送 HTTP 请求
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Sending request", i)
	}

	// 等待所有请求完成（简单示例，实际可用 WaitGroup）
	time.Sleep(10 * time.Second)
}
