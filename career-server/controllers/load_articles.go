package controllers

import (
	"career-server/lib/result"
	"career-server/models/page"

	"github.com/gin-gonic/gin"
)

func LoadArticles(ctx *gin.Context) {
	response := new(result.JsonResponseInterface)
	response.ErrNum = result.RESULT_SUCCESS
	response.ErrMsg = result.ErrInfos[response.ErrNum].ErrMsg
	page.LoadArticles(ctx, response)
	response.EchoResult(ctx)
}
