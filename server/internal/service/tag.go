package service

import (
	"sort"

	"github.com/huanglianjing/blog/server/internal/common"
	"github.com/huanglianjing/blog/server/internal/model"
)

// TagService 承载标签相关的业务逻辑。
type TagService struct {
	articleSvc *ArticleService
}

// NewTagService 构造 TagService。
func NewTagService() *TagService {
	return &TagService{articleSvc: NewArticleService()}
}

// TagOverviewResult 是标签概览接口的返回结构。
type TagOverviewResult struct {
	List []model.TagCount `json:"list"`
}

// Overview 返回所有标签及其文章数，按标签名排序：
// 数字 -> 字母 -> 汉字，字母不区分大小写，汉字按拼音升序。
func (s *TagService) Overview() (*TagOverviewResult, error) {
	list, err := model.ListTagsWithCount()
	if err != nil {
		return nil, err
	}
	keys := make(map[string]string, len(list))
	for _, t := range list {
		keys[t.Name] = common.NameSortKey(t.Name)
	}
	sort.SliceStable(list, func(i, j int) bool {
		return keys[list[i].Name] < keys[list[j].Name]
	})
	return &TagOverviewResult{List: list}, nil
}

// List 按标签名分页返回文章列表，日期倒序，每页 PageSize 条。
// 标签不存在时返回 (nil, nil)。
func (s *TagService) List(name string, page int) (*ArticleListResult, error) {
	if page < 0 {
		page = 0
	}

	tag, err := model.GetTagByName(name)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, nil
	}

	total, err := model.CountArticlesByTag(tag.ID)
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(PageSize) - 1) / int64(PageSize))

	articles, err := model.ListArticlesByTag(tag.ID, page*PageSize, PageSize)
	if err != nil {
		return nil, err
	}
	if err := s.articleSvc.enrichArticles(articles); err != nil {
		return nil, err
	}
	return &ArticleListResult{List: articles, TotalPages: totalPages}, nil
}
