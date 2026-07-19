package model

import (
	"errors"

	"gorm.io/gorm"
)

// Category 分类表
type Category struct {
	ID   int64  `json:"id"    gorm:"column:id;primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:name;type:text;not null;default:''"`
}

// TableName 表名
func (Category) TableName() string {
	return "category"
}

// CategoryCount 分类及其包含的文章数，用于分类概览。
type CategoryCount struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// ListCategoriesByIDs 按 id 批量查询分类。
func ListCategoriesByIDs(ids []int64) ([]Category, error) {
	categories := make([]Category, 0)
	if len(ids) == 0 {
		return categories, nil
	}
	if err := DB.Where("id IN ?", ids).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByName 按名称查询单个分类，未找到返回 (nil, nil)。
func GetCategoryByName(name string) (*Category, error) {
	var category Category
	err := DB.Where("name = ?", name).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// ListCategoriesWithCount 关联 article 表统计各分类的文章数，按文章数降序返回。
// 没有文章的分类不会出现在结果中。
func ListCategoriesWithCount() ([]CategoryCount, error) {
	result := make([]CategoryCount, 0)
	err := DB.Model(&Article{}).
		Select("category.name AS name, COUNT(article.id) AS count").
		Joins("JOIN category ON category.id = article.category_id").
		Group("category.id").
		Order("count DESC, category.name ASC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
