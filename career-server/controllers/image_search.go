package controllers

import (
	"career-server/lib/result"
	"career-server/models/page"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func ImageSearch(ctx *gin.Context) {
	response := new(result.JsonResponseInterface)
	response.ErrNum = result.RESULT_SUCCESS
	response.ErrMsg = result.ErrInfos[response.ErrNum].ErrMsg
	file, handler, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Fatal("ctx.Request.FormFile error,err=" + err.Error())
		response.ErrNum = result.RESULT_FILE_UPLOAD_ERROR
		response.ErrMsg = result.ErrInfos[response.ErrNum].ErrMsg
		response.EchoResult(ctx)
		return
	}
	fileName := handler.Filename
	curPath, _ := os.Getwd()
	// fmt.Println(curPath + "/images/" + fileName)
	loc, err := os.Create(curPath + "/images/" + fileName)
	if err != nil {
		log.Fatal("os.Create error,err=" + err.Error())
		response.ErrNum = result.RESULT_FILE_UPLOAD_ERROR
		response.ErrMsg = result.ErrInfos[response.ErrNum].ErrMsg
		response.EchoResult(ctx)
		return
	}
	_, err = io.Copy(loc, file)
	if err != nil {
		log.Fatal("image save error,err=" + err.Error())
		response.ErrNum = result.RESULT_FILE_UPLOAD_ERROR
		response.ErrMsg = result.ErrInfos[response.ErrNum].ErrMsg
		response.EchoResult(ctx)
		return
	}
	defer func() {
		err := loc.Close()
		if err != nil {
			log.Fatal("loc close error,err=" + err.Error())
		}
		// os.Remove(curPath + "/images/" + fileName)
	}()
	page.ImageSearch(ctx, curPath+"/images/"+fileName, response)
	response.EchoResult(ctx)
}
