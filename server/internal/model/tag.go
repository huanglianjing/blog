package model

import (
	"errors"

	"gorm.io/gorm"
)

// Tag 标签表
type Tag struct {
	ID   int64  `json:"id"    gorm:"column:id;primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:name;type:text;not null;default:''"`
}

// TableName 表名
func (Tag) TableName() string {
	return "tag"
}

// TagCount 标签及其包含的文章数，用于标签概览。
type TagCount struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// ListTagsByIDs 按 id 批量查询标签。
func ListTagsByIDs(ids []int64) ([]Tag, error) {
	tags := make([]Tag, 0)
	if len(ids) == 0 {
		return tags, nil
	}
	if err := DB.Where("id IN ?", ids).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// GetTagByName 按名称查询单个标签，未找到返回 (nil, nil)。
func GetTagByName(name string) (*Tag, error) {
	var tag Tag
	err := DB.Where("name = ?", name).First(&tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// ListTagsWithCount 关联 article_tag 表统计各标签的文章数。
// 没有文章的标签不会出现在结果中。排序由 service 层按名称规则处理。
func ListTagsWithCount() ([]TagCount, error) {
	result := make([]TagCount, 0)
	err := DB.Model(&Tag{}).
		Select("tag.name AS name, COUNT(article_tag.article_id) AS count").
		Joins("JOIN article_tag ON article_tag.tag_id = tag.id").
		Group("tag.id").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
