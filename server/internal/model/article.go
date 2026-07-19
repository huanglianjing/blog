package model

import (
	"errors"

	"gorm.io/gorm"
)

// Article 文章表
type Article struct {
	ID         int64  `json:"id"    gorm:"column:id;primaryKey;autoIncrement"`
	Title      string `json:"title" gorm:"column:title;type:text;not null;default:''"`
	Date       string `json:"date"  gorm:"column:date;type:text;not null;default:'';index"`
	Path       string `json:"path"  gorm:"column:path;type:text;not null;default:''"`
	CategoryID int64  `json:"category_id" gorm:"column:category_id;type:bigint;not null;default:0;index"`

	// 以下字段不是 article 表的列（gorm:"-"），
	// 而是查询时从 category / tag 表关联读取后填充。
	CategoryName string   `json:"category_name" gorm:"-"`
	Tags         []string `json:"tags"          gorm:"-"`
	Summary      string   `json:"summary"       gorm:"-"` // 正文纯文本开头预览
}

// TableName 表名
func (Article) TableName() string {
	return "article"
}

// ListArticles 按日期倒序分页查询 article 表，offset 起始偏移，limit 条数。
// 只查 article 单表，CategoryName / Tags 等关联字段由 service 层填充。
func ListArticles(offset, limit int) ([]Article, error) {
	articles := make([]Article, 0)
	if err := DB.Order("date DESC, id DESC").Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// CountArticles 返回 article 表的记录总数。
func CountArticles() (int64, error) {
	var count int64
	if err := DB.Model(&Article{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ListArticlesByCategory 按分类 id 分页查询文章，日期倒序。
func ListArticlesByCategory(categoryID int64, offset, limit int) ([]Article, error) {
	articles := make([]Article, 0)
	err := DB.Where("category_id = ?", categoryID).
		Order("date DESC, id DESC").Offset(offset).Limit(limit).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// CountArticlesByCategory 返回某分类下的文章总数。
func CountArticlesByCategory(categoryID int64) (int64, error) {
	var count int64
	if err := DB.Model(&Article{}).Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetArticleByTitle 按标题查询单篇文章，未找到返回 (nil, nil)。
func GetArticleByTitle(title string) (*Article, error) {
	var article Article
	err := DB.Where("title = ?", title).First(&article).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &article, nil
}
