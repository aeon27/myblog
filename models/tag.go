package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 用于GORM使用
type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 添加标签
func AddTag(name, createdBy string, state int) error {
	err := db.Create(&Tag{
		Name:      name,
		CreatedBy: createdBy,
		State:     state,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

// 分页获取标签
func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var tags []Tag
	var err error

	if pageNum > 0 && pageSize > 0 { // 分页查询
		err = db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	} else { // 全量查询
		err = db.Where(maps).Find(&tags).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

// 获取标签数量
func GetTagTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// 编辑标签
func EditTag(id int, data map[string]interface{}) error {
	err := db.Model(&Tag{}).Where("id = ?", id).Update(data).Error
	if err != nil {
		return err
	}

	return nil
}

// 删除标签
func DeleteTag(id int) error {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error
	if err != nil {
		return err
	}

	return nil
}

// 根据名字判断标签是否存在
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 根据id判断标签是否存在
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 硬删除tag，GORM约定硬删除用Unscoped
func CleanAllTags() error {
	err := db.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{}).Error
	if err != nil {
		return err
	}

	return nil
}

// 属于gorm的钩子机制
// 可以将回调方法定义为模型结构的指针，在创建、更新、查询、删除时将被调用
// 如果任何回调返回错误，gorm将停止未来操作并回滚所有更改
func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("CreatedOn", time.Now().Unix())

	return err
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("ModifiedOn", time.Now().Unix())

	return err
}
