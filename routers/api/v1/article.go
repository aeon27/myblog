package v1

import (
	"net/http"

	"github.com/aeon27/myblog/pkg/app"
	"github.com/aeon27/myblog/pkg/e"
	"github.com/aeon27/myblog/pkg/logging"
	"github.com/aeon27/myblog/pkg/setting"
	"github.com/aeon27/myblog/pkg/util"
	"github.com/aeon27/myblog/service/article_service"
	"github.com/aeon27/myblog/service/tag_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageURL string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// 添加文章
func AddArticle(c *gin.Context) {
	resp := &app.Responsor{GinContext: c}
	form := &AddArticleForm{}

	httpCode, errCode := app.BindAndValid(c, form)
	if errCode != e.SUCCESS {
		resp.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	tagExists, err := tagService.ExistByID()
	if err != nil {
		logging.Warn(err)
		resp.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_TAG_FAIL, nil)
		return
	}
	if !tagExists {
		resp.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CreatedBy:     form.CreatedBy,
		CoverImageUrl: form.CoverImageURL,
		State:         form.State,
	}
	if err := articleService.Add(); err != nil {
		resp.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	resp.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"MaxSize(100)"`
	Desc          string `form:"desc" valid:"MaxSize(255)"`
	Content       string `form:"content" valid:"MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageURL string `form:"cover_image_url" valid:"Requried;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// 编辑文章
func EditArticle(c *gin.Context) {
	resp := app.Responsor{GinContext: c}
	form := EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		resp.Response(httpCode, e.INVALID_PARAMS, nil)
		return
	}

	// 先校验tag是否存在
	tagService := tag_service.Tag{ID: form.TagID}
	tagExists, err := tagService.ExistByID()
	if err != nil {
		logging.Warn(err)
		resp.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_TAG_FAIL, nil)
		return
	}
	if !tagExists {
		resp.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		ModifiedBy:    form.ModifiedBy,
		CoverImageUrl: form.CoverImageURL,
		State:         form.State,
	}

	// 再校验对应文章是否存在
	articleExists, err := articleService.ExistByID()
	if err != nil {
		logging.Warn(err)
		resp.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !articleExists {
		resp.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		logging.Warn(err)
		resp.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	resp.Response(http.StatusOK, e.SUCCESS, nil)
}

// 获取文章
func GetArticle(c *gin.Context) {
	resp := app.Responsor{GinContext: c}

	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("id必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		resp.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		logging.Warn(err)
		resp.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		resp.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		logging.Warn(err)
		resp.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	resp.Response(http.StatusOK, e.SUCCESS, article)
}

// 获取文章列表
func GetArticles(c *gin.Context) {
	resp := app.Responsor{GinContext: c}
	valid := validation.Validation{}

	var state int = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("state只能为0或1")
	}

	var tagID int = -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagID = com.StrTo(arg).MustInt()
		valid.Min(tagID, 1, "tag_id").Message("tag_id必须大于0")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		resp.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		TagID:    tagID,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	count, err := articleService.GetCount()
	if err != nil {
		logging.Warn(err)
		resp.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_COUNT_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		logging.Warn(err)
		resp.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := map[string]interface{}{
		"lists": articles,
		"total": count,
	}

	resp.Response(http.StatusOK, e.SUCCESS, data)
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	resp := app.Responsor{GinContext: c}

	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("id必须大于0")

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		logging.Warn(err)
		resp.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	if !exists {
		resp.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		logging.Warn(err)
		resp.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
	}

	resp.Response(http.StatusOK, e.SUCCESS, nil)
}
