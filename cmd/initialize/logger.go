package initialize

import (
	"log"

	"github.com/xinliangnote/go-gin-api/cmd/global"
	"github.com/xinliangnote/go-gin-api/cmd/pkg/logger"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger(logFileName string) {
	global.Logger = NewLogger(logFileName)
}

func NewLogger(logFileName string) *logger.Logger {
	fileName := global.ServerConfig.AppInfo.LogSavePath + "/" + logFileName + global.ServerConfig.AppInfo.LogFileExt
	return logger.NewLogger(&lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   100,
		MaxAge:    30,
		LocalTime: true,
	}, "", log.LstdFlags)
}
