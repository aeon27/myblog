package app

import (
	"net/http"

	"github.com/aeon27/myblog/pkg/e"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// 根据结构体标签绑定数据到结构体并校验
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, e.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	return http.StatusOK, e.SUCCESS
}
