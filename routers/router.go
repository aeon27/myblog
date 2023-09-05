package routers

import (
	"net/http"

	"github.com/aeon27/myblog/middlewares/jwt"
	"github.com/aeon27/myblog/pkg/export"
	"github.com/aeon27/myblog/pkg/qrcode"
	"github.com/aeon27/myblog/pkg/setting"
	"github.com/aeon27/myblog/pkg/upload"
	"github.com/aeon27/myblog/routers/api"
	v1 "github.com/aeon27/myblog/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// gin模式设置
	gin.SetMode(setting.ServerSetting.RunMode)

	r := gin.New() // no middleware attatched
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 鉴权
	r.GET("/auth", api.GetAuth)

	// 图片上传和访问
	r.POST("/upload", api.UploadImage)                                // 上传图片
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath())) // 访问图片

	// 访问excel文件
	r.StaticFS("/export", http.Dir(export.GetExportFullPath()))

	// 访问二维码海报
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQRCodeFullPath()))

	apiv1 := r.Group("api/v1")
	// 注意，对非鉴权接口路由使用jwt中间件，此处即为api/v1
	apiv1.Use(jwt.JWT())
	{
		// 标签
		apiv1.GET("/tags", v1.GetTags)
		apiv1.POST("/tags", v1.AddTag)
		apiv1.PUT("/tags/:id", v1.EditTag)
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		r.POST("tags/export", v1.ExportTag) // 由数据库导出标签文件
		r.POST("tags/import", v1.ImportTag) // 由标签文件导入数据库

		// 文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		apiv1.GET("/articles", v1.GetArticles)
		apiv1.POST("/articles", v1.AddArticle)
		apiv1.PUT("/articles/:id", v1.EditArticle)
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
		apiv1.POST("/articles/poster/generate", v1.GenArticlePoster)
	}

	return r
}
