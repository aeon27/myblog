package tag_service

import (
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/aeon27/myblog/gredis"
	"github.com/aeon27/myblog/models"
	"github.com/aeon27/myblog/pkg/export"
	"github.com/aeon27/myblog/pkg/file"
	"github.com/aeon27/myblog/service/cache_service"
	"github.com/tealeg/xlsx"
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

// 将所有标签信息导出到文件
func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	xlsFile := xlsx.NewFile()
	sheet, err := xlsFile.AddSheet("标签信息")
	if err != nil {
		return "", nil
	}

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	for _, title := range titles { // 将标签属性作为列名加入第一行
		cell := row.AddCell()
		cell.Value = title
	}

	for _, tag := range tags { // 将每一个标签的所有属性信息作为新的row加入sheet
		values := []string{
			strconv.Itoa(tag.ID),
			tag.Name,
			tag.CreatedBy,
			strconv.Itoa(tag.CreatedOn),
			tag.ModifiedBy,
			strconv.Itoa(tag.ModifiedOn),
		}

		row := sheet.AddRow()
		for _, v := range values {
			cell := row.AddCell()
			cell.Value = v
		}
	}

	exportPath := export.GetExportFullPath()
	err = file.IsNotExistMkDir(exportPath) // 不存在则创建目录
	if err != nil {
		return "", err
	}

	time := strconv.Itoa(int(time.Now().Unix()))
	fileName := "tags-" + time + ".xlsx"
	err = xlsFile.Save(exportPath + fileName) // 导出到指定路径
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// 导入标签
func (t *Tag) Import(reader io.Reader) error {
	xlsx, err := excelize.OpenReader(reader)
	if err != nil {
		return err
	}

	rows := xlsx.GetRows("标签信息")
	for _, row := range rows {
		var data []string
		for _, cell := range row {
			data = append(data, cell)
		}
		if data[0] == "ID" {
			continue
		}
		err = models.AddTag(data[1], data[2], 1)
		if err != nil {
			return err
		}
	}

	return nil
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
