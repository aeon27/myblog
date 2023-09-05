package export

import "github.com/aeon27/myblog/pkg/setting"

// 获取导出excel完整的URL
func GetExportFullUrl(fileName string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExportPath() + fileName
}

// 获取用于导出excel的相对路径
func GetExportPath() string {
	return setting.AppSetting.ExportSavePath
}

// 获取用于导出excel的完整路径，即/runtime/export
func GetExportFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExportPath()
}
