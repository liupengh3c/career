package result

import (
	"sort"
)

// *****************************************************************//
// 新增错误码时，不仅在下面的常量组里增加错误常量						   //
// 还需要在ErrInfos字典中增加对应的错误明细，包括错误码、错误信息、错误说明   //
// *****************************************************************//
const (
	RESULT_SUCCESS              = 0
	RESULT_CREATE_BANNERS_ERROR = 1000
	RESULT_FILE_EMPTY_ERROR     = 1001
	RESULT_USER_NAME_ERROR      = 1002
	RESULT_USER_PASSWORD_ERROR  = 1003
	RESULT_SELECT_RECORD_ERROR  = 1004
	RESULT_USER_NOT_EXIST       = 1005
	RESULT_USER_ACCOUNT_ERROR   = 1006
	RESULT_USER_ENCRYPT_ERROR   = 1007
	RESULT_LOGIN_EXPIRE_ERROR   = 1008
	RESULT_REDIS_WRITE_ERROR    = 1009
	RESULT_USER_NOTLOGIN_ERROR  = 1010
	RESULT_INTERNAL_ERROR       = 1011
	RESULT_FILE_UPLOAD_ERROR    = 1012
	RESULT_PARAM_ERROR          = 1013
	RESULT_NOT_AUTHOTITY_ERROR  = 1014
	RESULT_GET_ORDER_ERROR      = 1015
)

type ErrInfo struct {
	ErrNum int    `json:"err_num"`
	ErrMsg string `json:"err_msg"`
	Help   string `json:"help_info"`
}
type errSlice []ErrInfo

// 必须要在这里增加错误码对应的错误明细
var ErrInfos = map[int]ErrInfo{
	0: {
		ErrNum: 0,
		ErrMsg: "操作成功",
		Help:   "操作成功",
	},
	1000: {
		ErrNum: 1000,
		ErrMsg: "创建banners失败",
		Help:   "写数据库失败",
	},
	1001: {
		ErrNum: 1001,
		ErrMsg: "文件内容不能为空",
		Help:   "模拟量报警文件,内容为空",
	},
	1002: {
		ErrNum: 1002,
		ErrMsg: "用户名不能为空",
		Help:   "用户名不能为空",
	},
	1003: {
		ErrNum: 1003,
		ErrMsg: "密码不能为空",
		Help:   "密码不能为空",
	},
	1004: {
		ErrNum: 1004,
		ErrMsg: "账号登录失败,请重试.",
		Help:   "读取数据库失败",
	},
	1005: {
		ErrNum: 1005,
		ErrMsg: "账号不存在,请重新输入.",
		Help:   "账号不存在,数据库中没有该账号",
	},
	1006: {
		ErrNum: 1006,
		ErrMsg: "用户名或密码错误,请重新输入.",
		Help:   "密码错误",
	},
	1007: {
		ErrNum: 1007,
		ErrMsg: "账号登录失败,请重试.",
		Help:   "账号密码加密失败.",
	},
	1008: {
		ErrNum: 1008,
		ErrMsg: "登录已失效,请重新登录.",
		Help:   "登录已失效,请重新登录.",
	},
	1009: {
		ErrNum: 1009,
		ErrMsg: "账号登录失败,请重试.",
		Help:   "写redis失败.",
	},
	1010: {
		ErrNum: 1010,
		ErrMsg: "账号未登录,请登录.",
		Help:   "账号未登录,请登录.",
	},
	1011: {
		ErrNum: 1011,
		ErrMsg: "网络不给力,请稍后重试.",
		Help:   "网络不给力,请稍后重试.",
	},
	1012: {
		ErrNum: 1012,
		ErrMsg: "文件上传失败,请稍后重试.",
		Help:   "文件上传失败,请稍后重试.",
	},
	1013: {
		ErrNum: 1013,
		ErrMsg: "参数错误.",
		Help:   "缺少参数或者为空.",
	},
	1014: {
		ErrNum: 1014,
		ErrMsg: "没有权限，请联系管理员开通.",
		Help:   "没有权限，请联系管理员开通.",
	},
	1015: {
		ErrNum: 1015,
		ErrMsg: "获取订单列表失败.",
		Help:   "获取订单列表失败.",
	},
}

func (Info errSlice) Len() int {
	return len(Info)
}
func (Info errSlice) Swap(i, j int) {
	Info[i], Info[j] = Info[j], Info[i]
}
func (Info errSlice) Less(i, j int) bool {
	return Info[i].ErrNum < Info[j].ErrNum
}

/*
pn:页码
rn：每页获取条数
*/
func GetErrCodes(pn, rn int) (int, errSlice) {
	var errCodes, errTmp errSlice
	var start, end int
	// 先获取所有错误码并按错误码从小到大排序
	for _, value := range ErrInfos {
		errTmp = append(errTmp, value)
	}
	sort.Stable(errTmp)
	// 根据pn、rn获取对应的错误码，需要进行越界保护
	if rn == 0 {
		start = 0
		end = len(errTmp)
	} else {
		start = (pn - 1) * rn
		end = start + rn
		if len(errTmp) <= start {
			return RESULT_SUCCESS, errCodes
		}
		if len(errTmp) <= end {
			end = len(errTmp)
		}
	}
	for _, value := range errTmp[start:end] {
		errCodes = append(errCodes, value)
	}
	return RESULT_SUCCESS, errCodes
}
