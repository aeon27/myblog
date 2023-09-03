package app

import (
	"github.com/aeon27/myblog/pkg/logging"
	"github.com/astaxie/beego/validation"
)

// 将参数校验发生的错误记录在本地日志
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}
}
