package model

import (
	"errors"

	"gorm.io/gorm"

	"github.com/huanglianjing/blog/server/internal/common"
)

// Article 文章表
type Article struct {
	ID         int64  `json:"id"    gorm:"column:id;primaryKey;autoIncrement"`
	Title      string `json:"title" gorm:"column:title;type:text;not null;default:''"`
	Date       string `json:"date"  gorm:"column:date;type:text;not null;default:'';index"`
	Path       string `json:"path"  gorm:"column:path;type:text;not null;default:''"`
	CategoryID int64  `json:"category_id" gorm:"column:category_id;type:bigint;not null;default:0;index"`

	// Content 是正文纯文本，仅供搜索匹配与片段截取，不返回给前端。
	// 由 article_converter 写入。该列体积大，不需要正文的查询应 Omit("content")。
	// 不建索引：子串 LIKE 用不上索引。
	Content string `json:"-" gorm:"column:content;type:text;not null;default:''"`

	// 以下字段不是 article 表的列（gorm:"-"），
	// 而是查询时从 category / tag 表关联读取后填充。
	CategoryName string   `json:"category_name" gorm:"-"`
	Tags         []string `json:"tags"          gorm:"-"`
	Summary      string   `json:"summary"       gorm:"-"`     // 正文纯文本开头预览
	Snippet      string   `json:"snippet,omitempty" gorm:"-"` // 搜索命中处的上下文片段
}

// TableName 表名
func (Article) TableName() string {
	return "article"
}

// ListArticles 按日期倒序分页查询 article 表，offset 起始偏移，limit 条数。
// 只查 article 单表，CategoryName / Tags 等关联字段由 service 层填充。
func ListArticles(offset, limit int) ([]Article, error) {
	articles := make([]Article, 0)
	if err := DB.Omit("content").Order("date DESC, id DESC").Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
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
	err := DB.Omit("content").Where("category_id = ?", categoryID).
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
	err := DB.Omit("content").Where("title = ?", title).First(&article).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// SearchArticlesByTitle 按标题子串匹配分页查询文章，日期倒序。
func SearchArticlesByTitle(keyword string, offset, limit int) ([]Article, error) {
	articles := make([]Article, 0)
	err := DB.Omit("content").Where(`title LIKE ? ESCAPE '\'`, common.LikePattern(keyword)).
		Order("date DESC, id DESC").Offset(offset).Limit(limit).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// CountArticlesByTitle 返回标题子串匹配的文章总数。
func CountArticlesByTitle(keyword string) (int64, error) {
	var count int64
	err := DB.Model(&Article{}).Where(`title LIKE ? ESCAPE '\'`, common.LikePattern(keyword)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SearchArticlesByContent 按正文子串匹配分页查询文章，日期倒序。
// 排除标题已命中的文章，避免与标题类结果重复；
// 不 Omit content，service 层需要正文来截取命中片段。
func SearchArticlesByContent(keyword string, offset, limit int) ([]Article, error) {
	articles := make([]Article, 0)
	pattern := common.LikePattern(keyword)
	err := DB.Where(`content LIKE ? ESCAPE '\' AND title NOT LIKE ? ESCAPE '\'`, pattern, pattern).
		Order("date DESC, id DESC").Offset(offset).Limit(limit).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// CountArticlesByContent 返回正文子串匹配（且标题未命中）的文章总数。
func CountArticlesByContent(keyword string) (int64, error) {
	var count int64
	pattern := common.LikePattern(keyword)
	err := DB.Model(&Article{}).
		Where(`content LIKE ? ESCAPE '\' AND title NOT LIKE ? ESCAPE '\'`, pattern, pattern).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountArticlesWithContent 返回正文纯文本非空的文章数，
// 用于服务启动时检查数据库是否已由新版 article_converter 写入正文。
func CountArticlesWithContent() (int64, error) {
	var count int64
	if err := DB.Model(&Article{}).Where("content != ''").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
