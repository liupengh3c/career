package logservice

import (
	"career-server/lib/tools"
	"runtime"
	"strconv"

	"github.com/cihub/seelog"
)

var logger seelog.LoggerInterface

func getLogPrefix(level int) string {
	pc, file, line, _ := runtime.Caller(level)
	info := "[" + file + ":" + strconv.Itoa(line) + "] "
	function := runtime.FuncForPC(pc).Name()
	info += "[" + function + "]"
	return info
}
func Notice(strMsg string) {
	logger = seelog.Disabled
	path := tools.GetCurrentDirectory()
	logger, err := seelog.LoggerFromConfigAsFile(path + "/conf/log.xml")
	if err != nil {
		return
	}
	prefix := getLogPrefix(2)
	seelog.ReplaceLogger(logger)
	seelog.Trace(prefix + strMsg)
	defer seelog.Flush()
}

func Fatal(strMsg string) {
	logger = seelog.Disabled
	path := tools.GetCurrentDirectory()
	logger, err := seelog.LoggerFromConfigAsFile(path + "/conf/log.xml")
	if err != nil {
		return
	}
	seelog.ReplaceLogger(logger)
	prefix := getLogPrefix(2)
	seelog.Error(prefix + strMsg)
	defer seelog.Flush()
}

func Debug(strMsg string) {
	logger = seelog.Disabled
	path := tools.GetCurrentDirectory()
	logger, err := seelog.LoggerFromConfigAsFile(path + "/conf/log.xml")
	if err != nil {
		return
	}
	seelog.ReplaceLogger(logger)
	prefix := getLogPrefix(2)
	seelog.Debug(prefix + strMsg)
	defer seelog.Flush()
}
