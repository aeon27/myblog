package api

import (
	"net/http"

	"github.com/aeon27/myblog/pkg/app"
	"github.com/aeon27/myblog/pkg/e"
	"github.com/aeon27/myblog/pkg/logging"
	"github.com/aeon27/myblog/pkg/upload"
	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	resp := app.Responsor{GinContext: c}
	code := e.SUCCESS
	data := make(map[string]interface{})

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		resp.Response(http.StatusOK, code, data)
		return
	}

	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			if err := upload.CheckImage(fullPath); err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullURL(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	resp.Response(http.StatusOK, code, data)
}
