package article_service

import (
	"encoding/json"

	"github.com/aeon27/myblog/gredis"
	"github.com/aeon27/myblog/models"
	"github.com/aeon27/myblog/pkg/logging"
	"github.com/aeon27/myblog/service/cache_service"
)

type Article struct {
	ID    int
	TagID int

	State         int
	Title         string
	Desc          string
	Content       string
	CreatedBy     string
	ModifiedBy    string
	CoverImageUrl string

	PageNum  int
	PageSize int
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"state":           a.State,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
	}

	if err := models.AddArticle(article); err != nil {
		return err
	}

	return nil
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	// 先查redis缓存
	cacheService := cache_service.Article{ID: a.ID}
	key := cacheService.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil { // 缓存没有就记录本地日志，然后去数据库查
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	// 缓存没有去数据库查
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	// 刷新redis
	gredis.Set(key, article, 3600)

	return article, nil
}

func (a *Article) GetAll() ([]models.Article, error) {
	var articles, cacheArticles []models.Article

	cacheService := cache_service.Article{
		TagID: a.TagID,
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}

	// 先查缓存
	key := cacheService.GetArticlesKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	// 再查数据库
	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	// 刷新缓存
	gredis.Set(key, articles, 3600)

	return articles, nil
}

func (a *Article) Edit() error {
	maps := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"modified_by":     a.ModifiedBy,
		"state":           a.State,
	}
	return models.EditArticle(a.ID, maps)
}

func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

func (a *Article) GetCount() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistedArticleByID(a.ID)
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	if a.State >= 0 {
		maps["state"] = a.State
	}

	return maps
}
