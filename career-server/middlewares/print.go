package middlewares

import (
	"career-server/logservice"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
