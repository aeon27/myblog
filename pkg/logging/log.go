package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2 // 默认调用深度

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func init() {
	filePath := getLogFileFullPath()
	F = openLogFile(filePath)

	logger = log.New(F, DefaultPrefix, log.LstdFlags) // log.LstdFlags 标准日志格式 日期+时间
}

func Debug(v ...interface{}) {
	log.Println(v...)
	setPrefix(DEBUG)
	logger.Println(v...)
}

func Info(v ...interface{}) {
	log.Println(v...)
	setPrefix(INFO)
	logger.Println(v...)
}

func Warn(v ...interface{}) {
	log.Println(v...)
	setPrefix(WARN)
	logger.Println(v...)
}

func Error(v ...interface{}) {
	log.Println(v...)
	setPrefix(ERROR)
	logger.Println(v...)
}

func Fatal(v ...interface{}) {
	log.Println(v...)
	setPrefix(FATAL)
	logger.Println(v...)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], file, line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
