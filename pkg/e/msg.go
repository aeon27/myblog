package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	ERROR_EXIST_TAG:            "已存在该标签名称",
	ERROR_CHECK_EXIST_TAG_FAIL: "校验标签是否已存在失败",
	ERROR_NOT_EXIST_TAG:        "该标签不存在",
	ERROR_ADD_TAG_FAIL:         "添加标签失败",
	ERROR_EDIT_TAG_FAIL:        "编辑标签失败",
	ERROR_GET_TAGS_FAIL:        "获取所有标签失败",
	ERROR_GET_TAG_COUNT_FAIL:   "获取标签数量失败",
	ERROR_DELETE_TAG_FAIL:      "删除标签失败",
	ERROR_EXPORT_TAG_FAIL:      "导出标签失败",
	ERROR_TAG_ALREADY_EXISTS:   "该标签名已存在",
	ERROR_GET_TAG_FILE:         "获取标签文件失败",
	ERROR_IMPORT_TAG_FAIL:      "导入标签失败",

	ERROR_ADD_ARTICLE_FAIL:         "添加文章失败",
	ERROR_GET_ARTICLE_FAIL:         "获取文章失败",
	ERROR_NOT_EXIST_ARTICLE:        "该文章不存在",
	ERROR_GET_ARTICLES_FAIL:        "批量获取文章失败",
	ERROR_EDIT_ARTICLE_FAIL:        "编辑文章失败",
	ERROR_DELETE_ARTICLE_FAIL:      "删除文章失败",
	ERROR_GET_ARTICLE_COUNT_FAIL:   "获取文章数量失败",
	ERROR_CHECK_EXIST_ARTICLE_FAIL: "校验文章是否已存在失败",
	ERROR_GEN_ARTICLE_POSTER_FAIL:  "生成文章海报失败",

	ERROR_AUTH_CHECK_FAIL:          "鉴权参数校验失败",
	ERROR_AUTH_NOT_HAVE_TOKEN:      "Token缺失",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token鉴权超时",
	ERROR_AUTH_GEN_TOKEN_FAIL:      "Token生成失败",

	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",

	ERROR_GEN_QRCODE_FAIL: "生成二维码失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[ERROR]
	}
	return msg
}
