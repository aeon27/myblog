package tag_service

import (
	"encoding/json"

	"github.com/aeon27/myblog/gredis"
	"github.com/aeon27/myblog/models"
	"github.com/aeon27/myblog/service/cache_service"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.CreatedBy, t.State)
}

func (t *Tag) Edit() error {
	data := map[string]interface{}{}
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	data["state"] = t.State

	return models.EditTag(t.ID, data)
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var tags, cacheTags []models.Tag

	cacheService := cache_service.Tag{
		Name:  t.Name,
		State: t.State,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}

	// 判断redis缓存有无数据，从缓存获取tags
	key := cacheService.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			return nil, err
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	// 从数据库获取tags
	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	// 刷新缓存
	gredis.Set(key, tags, 3600)

	return tags, nil
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) GetCount() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}
