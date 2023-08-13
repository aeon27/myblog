package jwt

import (
	"net/http"
	"time"

	"github.com/aeon27/myblog/pkg/e"
	"github.com/aeon27/myblog/pkg/util"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = e.ERROR_AUTH_NOT_HAVE_TOKEN
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"data": data,
				"msg":  e.GetMsg(code),
			})

			// 忘记以下步骤，导致返回了2次c.JSON()
			// 鉴权不通过，调用c.Abort()中断后续中间件调用
			c.Abort()
			return
		}

		c.Next()
	}
}
