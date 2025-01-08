package result

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 此结构用于返回二维数组，data为二维切片。每一维为字典
type JsonResponseSliceMap struct {
	ErrNum int                 `json:"err_num"`
	ErrMsg string              `json:"err_msg"`
	Data   []map[string]string `json:"data"`
}

// 返回字典
type JsonResponseMap struct {
	ErrNum int               `json:"err_num"`
	ErrMsg string            `json:"err_msg"`
	Data   map[string]string `json:"data"`
}

// 返回字符串
type JsonResponseString struct {
	ErrNum int    `json:"err_num"`
	ErrMsg string `json:"err_msg"`
	Data   string `json:"data"`
}

// 混合类型
type JsonResponseInterface struct {
	ErrNum int         `json:"err_num"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}

// 混合类型
type JsonSliceResponseInterface struct {
	ErrNum int           `json:"err_num"`
	ErrMsg string        `json:"err_msg"`
	Data   []interface{} `json:"data"`
}

/*
code：错误码
ctx:webcontext，用于输出信息
*/
func (response *JsonSliceResponseInterface) EchoResult(ctx *gin.Context) {
	// response.ErrNum = ErrInfos[code].ErrNum
	// response.ErrMsg = ErrInfos[code].ErrMsg
	ctx.JSON(http.StatusOK, *response)
}

func (response *JsonSliceResponseInterface) Init() {
	response.ErrNum = RESULT_SUCCESS
	response.ErrMsg = ErrInfos[response.ErrNum].ErrMsg
	response.Data = make([]interface{}, 1)
}
func (response *JsonSliceResponseInterface) SetErrAndMsg(ctx *gin.Context, code int) {
	response.ErrNum = ErrInfos[code].ErrNum
	response.ErrMsg = ErrInfos[code].ErrMsg
	// ctx.JSON(http.StatusOK, *response)
}

/*
code：错误码
ctx:webcontext，用于输出信息
*/
func (response *JsonResponseSliceMap) EchoResult(ctx *gin.Context) {
	// response.ErrNum = ErrInfos[code].ErrNum
	// response.ErrMsg = ErrInfos[code].ErrMsg
	ctx.JSON(http.StatusOK, *response)
}

func (response *JsonResponseSliceMap) Init() {
	response.ErrNum = RESULT_SUCCESS
	response.ErrMsg = ErrInfos[response.ErrNum].ErrMsg
	response.Data = make([]map[string]string, 1)
}
func (response *JsonResponseSliceMap) SetErrInfos(ctx *gin.Context, code int) {
	response.ErrNum = ErrInfos[code].ErrNum
	response.ErrMsg = ErrInfos[code].ErrMsg
	ctx.JSON(http.StatusOK, *response)
}

/*
code：错误码
ctx:webcontext，用于输出信息
*/
func (response *JsonResponseMap) EchoResult(ctx *gin.Context) {
	// response.ErrNum = ErrInfos[code].ErrNum
	// response.ErrMsg = ErrInfos[code].ErrMsg
	ctx.JSON(http.StatusOK, *response)
}
func (response *JsonResponseMap) SetErrAndMsg(ctx *gin.Context, code int) {
	response.ErrNum = ErrInfos[code].ErrNum
	response.ErrMsg = ErrInfos[code].ErrMsg
	ctx.JSON(http.StatusOK, *response)
}
func (response *JsonResponseMap) Init() {
	response.ErrNum = RESULT_SUCCESS
	response.ErrMsg = ErrInfos[response.ErrNum].ErrMsg
	response.Data = make(map[string]string)
}

/*
code：错误码
ctx:webcontext，用于输出信息
*/
func (response *JsonResponseString) EchoResult(ctx *gin.Context) {
	// response.ErrNum = ErrInfos[code].ErrNum
	// response.ErrMsg = ErrInfos[code].ErrMsg
	ctx.JSON(http.StatusOK, *response)
}
func (response *JsonResponseString) Init() {
	response.ErrNum = RESULT_SUCCESS
	response.ErrMsg = ErrInfos[response.ErrNum].ErrMsg
	response.Data = ""
}

/*
code：错误码
ctx:webcontext，用于输出信息
*/
func (response *JsonResponseInterface) EchoResult(ctx *gin.Context) {
	// response.ErrNum = ErrInfos[code].ErrNum
	// response.ErrMsg = ErrInfos[code].ErrMsg
	ctx.JSON(http.StatusOK, *response)
}
func (response *JsonResponseInterface) Init() {
	response.ErrNum = RESULT_SUCCESS
	response.ErrMsg = ErrInfos[response.ErrNum].ErrMsg
}
func (response *JsonResponseInterface) SetErrAndMsg(ctx *gin.Context, code int) {
	response.ErrNum = ErrInfos[code].ErrNum
	response.ErrMsg = ErrInfos[code].ErrMsg
}

type JSONResponse struct {
	ErrNum int         `json:"err_num"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}

func NewJSONResponse(errno int, msg string, d interface{}) *JSONResponse {
	return &JSONResponse{
		ErrNum: errno,
		ErrMsg: msg,
		Data:   d,
	}
}

// 通过jsonp的例子证明以后可以通过JSONResponse结构体来适配不同的返回类型
func (r *JSONResponse) JsonpResponse(callback string) map[string]interface{} {
	d, _ := json.Marshal(r.Data)
	m := map[string]interface{}{
		"err_num": r.ErrNum,
		"err_msg": r.ErrMsg,
		"data":    callback + "(" + string(d) + ")",
	}
	return m
}
