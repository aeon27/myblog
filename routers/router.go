package router

import (
	"github.com/aeon27/myblog/middlewares/jwt"
	"github.com/aeon27/myblog/pkg/setting"
	"github.com/aeon27/myblog/routers/api"
	v1 "github.com/aeon27/myblog/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New() // no middleware attatched
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.Use(jwt.JWT())
	gin.SetMode(setting.ServerSetting.RunMode)

	// 鉴权
	r.GET("/auth", api.GetAuth)

	apiv1 := r.Group("api/v1")
	// 注意，对非鉴权接口路由使用jwt中间件，此处即为api/v1
	apiv1.Use(jwt.JWT())
	{
		// 标签
		apiv1.GET("/tags", v1.GetTags)
		apiv1.POST("/tags", v1.AddTag)
		apiv1.PUT("/tags/:id", v1.EditTag)
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		// 文章
		apiv1.GET("/articles", v1.GetArticles)
		apiv1.GET("/articles/:id", v1.GetArticle)
		apiv1.POST("/articles", v1.AddArticle)
		apiv1.PUT("/articles/:id", v1.EditArticle)
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
