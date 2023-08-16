package logging

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aeon27/myblog/pkg/setting"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", setting.AppSetting.LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt)

	return fmt.Sprintf("%s/%s", prefixPath, suffixPath)
}

func openLogFile(filePath string) *os.File {
	// os.Stat：返回文件信息结构描述文件。如果出现错误，会返回*PathError
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err): // 不存在则创建目录
		mkDir()
	case os.IsPermission(err): // 权限不满足
		log.Fatalf("Permission: %v", err)
	}

	// 0644 = -rw-r--r--; 4-r, 2-w, 1-x
	fileHandle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile: %v", err)
	}

	return fileHandle
}

func mkDir() {
	// os.Getwd()返回与当前目录对应的根路径名
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm) // ModePerm = 0777
	if err != nil {
		panic(err)
	}
}
