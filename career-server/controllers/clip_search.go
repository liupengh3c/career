package controllers

import (
	"career-server/lib/result"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func ClipSearch(ctx *gin.Context) {
	response := new(result.JsonResponseInterface)
	response.ErrNum = result.RESULT_SUCCESS
	response.ErrMsg = result.ErrInfos[response.ErrNum].ErrMsg
	body, _ := io.ReadAll(ctx.Request.Body)
	fmt.Println(string(body))
	response.EchoResult(ctx)
}
