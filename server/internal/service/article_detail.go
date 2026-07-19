package service

import (
	"fmt"
	"os"

	"github.com/huanglianjing/blog/server/internal/model"
)

// ArticleDetailResult 是文章详情接口的返回结构。
type ArticleDetailResult struct {
	Title        string   `json:"title"`
	Date         string   `json:"date"`          // 年月日，形如 2026-05-05
	CategoryName string   `json:"category_name"` // 所属分类
	Tags         []string `json:"tags"`          // 所属标签
	Content      string   `json:"content"`       // 正文 html
}

// Detail 按标题返回文章详情：元信息 + 从 path 读取的 html 正文。
// 文章不存在时返回 (nil, nil)。
func (s *ArticleService) Detail(title string) (*ArticleDetailResult, error) {
	article, err := model.GetArticleByTitle(title)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, nil
	}

	// 读取 path 指向的 html 正文。
	content, err := os.ReadFile(article.Path)
	if err != nil {
		return nil, fmt.Errorf("read html %q: %w", article.Path, err)
	}

	// 分类名。
	categoryNames, err := s.categoryNameMap([]int64{article.CategoryID})
	if err != nil {
		return nil, err
	}

	// 标签。
	articleTags, err := s.articleTagsMap([]int64{article.ID})
	if err != nil {
		return nil, err
	}
	tags := articleTags[article.ID]
	if tags == nil {
		tags = []string{}
	}

	return &ArticleDetailResult{
		Title:        article.Title,
		Date:         dateOnly(article.Date),
		CategoryName: categoryNames[article.CategoryID],
		Tags:         tags,
		Content:      string(content),
	}, nil
}

// dateOnly 从 "2026-05-05 19:23:46" 这样的日期时间中取出年月日 "2026-05-05"。
// 长度不足时原样返回。
func dateOnly(date string) string {
	if len(date) >= 10 {
		return date[:10]
	}
	return date
}
