package app

import (
	"github.com/aeon27/myblog/pkg/e"
	"github.com/gin-gonic/gin"
)

type Responsor struct {
	GinContext *gin.Context
}

func (resp *Responsor) Response(httpCode, errCode int, data interface{}) {
	resp.GinContext.JSON(httpCode, gin.H{
		"code": errCode,
		"data": data,
		"msg":  e.GetMsg(errCode),
	})
}
