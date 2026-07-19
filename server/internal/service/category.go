package service

import (
	"github.com/huanglianjing/blog/server/internal/model"
)

// CategoryService 承载分类相关的业务逻辑。
type CategoryService struct {
	articleSvc *ArticleService
}

// NewCategoryService 构造 CategoryService。
func NewCategoryService() *CategoryService {
	return &CategoryService{articleSvc: NewArticleService()}
}

// CategoryOverviewResult 是分类概览接口的返回结构。
type CategoryOverviewResult struct {
	List []model.CategoryCount `json:"list"`
}

// Overview 返回所有分类及其文章数，按文章数降序排列。
func (s *CategoryService) Overview() (*CategoryOverviewResult, error) {
	list, err := model.ListCategoriesWithCount()
	if err != nil {
		return nil, err
	}
	return &CategoryOverviewResult{List: list}, nil
}

// List 按分类名分页返回文章列表，日期倒序，每页 PageSize 条。
// 分类不存在时返回 (nil, nil)。
func (s *CategoryService) List(name string, page int) (*ArticleListResult, error) {
	if page < 0 {
		page = 0
	}

	category, err := model.GetCategoryByName(name)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, nil
	}

	total, err := model.CountArticlesByCategory(category.ID)
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(PageSize) - 1) / int64(PageSize))

	articles, err := model.ListArticlesByCategory(category.ID, page*PageSize, PageSize)
	if err != nil {
		return nil, err
	}
	if err := s.articleSvc.enrichArticles(articles); err != nil {
		return nil, err
	}
	return &ArticleListResult{List: articles, TotalPages: totalPages}, nil
}
