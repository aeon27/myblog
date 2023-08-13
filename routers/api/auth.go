package api

import (
	"net/http"

	"github.com/aeon27/myblog/models"
	"github.com/aeon27/myblog/pkg/e"
	"github.com/aeon27/myblog/pkg/logging"
	"github.com/aeon27/myblog/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	if ok, _ := valid.Valid(&a); ok {
		if models.CheckAuth(username, password) {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_GEN_TOKEN_FAIL
			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH_CHECK_FAIL
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  e.GetMsg(code),
	})
}
