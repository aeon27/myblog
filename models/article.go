package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"` // gorm:"index" 用于声明这个字段为索引，如果使用了自动迁移功能则会有所影响
	Tag   Tag `json:"tag"`

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 根据文章id判断是否存在
func ExistedArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(&article)
	if article.ID > 0 {
		return true
	}

	return false
}

// 添加文章
func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	})

	return true
}

// 获取文章
func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	// Article有一个结构体成员是TagID，就是外键。
	// gorm会通过 类名+ID 的方式去找到这两个类之间的关联关系
	// Article有一个嵌套在里的Tag结构体，我们可以通过Related进行关联查询
	db.Model(&article).Related(&article.Tag)

	return
}

// 获取文章列表
func GetArticles(pageNum, pageSize int, maps interface{}) (articles []Article) {
	// Preload就是一个预加载器，它会执行两条 SQL
	// 分别是SELECT * FROM blog_articles;
	// 和SELECT * FROM blog_tag WHERE id IN (1,2,3,4);
	// 那么在查询出结构后，gorm内部处理对应的映射逻辑，将其填充到Article的Tag中，会特别方便，并且避免了循环查询
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

// 获取文章数量
func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)

	return
}

// 编辑文章
func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Update(data)

	return true
}

// 删除文章
func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(&Article{})

	return true
}

// 属于gorm的钩子机制
// 可以将回调方法定义为模型结构的指针，在创建、更新、查询、删除时将被调用
// 如果任何回调返回错误，gorm将停止未来操作并回滚所有更改
func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("CreatedOn", time.Now().Unix())

	return err
}

func (article *Article) AfterCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("ModifiedOn", time.Now().Unix())

	return err
}
