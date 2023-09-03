package v1

import (
	"net/http"

	"github.com/aeon27/myblog/pkg/app"
	"github.com/aeon27/myblog/pkg/e"
	"github.com/aeon27/myblog/pkg/setting"
	"github.com/aeon27/myblog/pkg/util"
	"github.com/aeon27/myblog/service/tag_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func GetTags(c *gin.Context) {
	resp := app.Responsor{GinContext: c}

	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	tags, err := tagService.GetAll()
	if err != nil {
		resp.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.GetCount()
	if err != nil {
		resp.Response(http.StatusInternalServerError, e.ERROR_GET_TAG_COUNT_FAIL, nil)
		return
	}

	data := map[string]interface{}{
		"lists": tags,
		"total": count,
	}

	resp.Response(http.StatusOK, e.SUCCESS, data)
}

type AddTagForm struct {
	Name       string `json:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `json:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `json:"state" valid:"Range(0,1)"`
}

func AddTag(c *gin.Context) {
	resp := app.Responsor{GinContext: c}
	form := EditTagForm{}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		resp.Response(httpCode, errCode, nil)
		return
	}

	tagService := &tag_service.Tag{
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistByName()
	if err != nil {
		resp.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		resp.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		resp.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	resp.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `json:"id" valid:"Required;Min(1)"`
	Name       string `json:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `json:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `json:"state" valid:"Range(0,1)"`
}

func EditTag(c *gin.Context) {
	resp := app.Responsor{GinContext: c}
	form := EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()} // 注意Query和Param的区别

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		resp.Response(httpCode, errCode, nil)
		return
	}

	tagService := &tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		resp.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		resp.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		resp.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	resp.Response(http.StatusOK, e.SUCCESS, nil)
}

func DeleteTag(c *gin.Context) {
	resp := app.Responsor{GinContext: c}

	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		resp.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		ID: id,
	}
	exists, err := tagService.ExistByID()
	if err != nil {
		resp.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		resp.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Delete()
	if err != nil {
		resp.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	resp.Response(http.StatusOK, e.SUCCESS, nil)
}
