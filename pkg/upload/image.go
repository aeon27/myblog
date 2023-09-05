package upload

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aeon27/myblog/pkg/file"
	"github.com/aeon27/myblog/pkg/logging"
	"github.com/aeon27/myblog/pkg/setting"
	"github.com/aeon27/myblog/pkg/util"
)

// 获取图片完整的URL
func GetImageFullURL(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

// 获取图片名称，将文件名部分转为MD5，后缀不变
func GetImageName(name string) string {
	ext := file.GetExt(name)
	fileName := util.EncodeMD5(strings.TrimSuffix(name, ext))

	return fileName + ext
}

// 获取图片用于保存的相对路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

// 获取图片用于保存的完整路径
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

// 检查图片扩展名
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

// 检查图片大小
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

// 检查文件目录权限，没有则创建
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	noPerm := file.CheckPermission(src)
	if noPerm {
		return fmt.Errorf("file.CheckPermission Permission denied src: %v", err)
	}

	return nil
}
