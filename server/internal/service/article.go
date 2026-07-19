package service

import (
	"os"

	"github.com/huanglianjing/blog/server/internal/common"
	"github.com/huanglianjing/blog/server/internal/model"
)

// ArticleService 承载文章相关的业务逻辑。
type ArticleService struct{}

// NewArticleService 构造 ArticleService。
func NewArticleService() *ArticleService {
	return &ArticleService{}
}

// PageSize 是文章列表每页固定的记录数。
const PageSize = 10

// ArticleListResult 是 List 接口的返回结构。
type ArticleListResult struct {
	List       []model.Article `json:"list"`
	TotalPages int             `json:"total_pages"`
}

// List 按页返回文章记录（page 从 0 开始，每页 PageSize 条），
// 并填充分类名与标签等关联字段，同时返回总页数。
func (s *ArticleService) List(page int) (*ArticleListResult, error) {
	if page < 0 {
		page = 0
	}

	total, err := model.CountArticles()
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(PageSize) - 1) / int64(PageSize))

	articles, err := model.ListArticles(page*PageSize, PageSize)
	if err != nil {
		return nil, err
	}
	if err := s.enrichArticles(articles); err != nil {
		return nil, err
	}
	return &ArticleListResult{List: articles, TotalPages: totalPages}, nil
}

// enrichArticles 就地填充文章的日期、分类名、标签与预览等关联字段。
func (s *ArticleService) enrichArticles(articles []model.Article) error {
	if len(articles) == 0 {
		return nil
	}

	// 收集用到的文章 id 和分类 id。
	articleIDs := make([]int64, 0, len(articles))
	categoryIDs := make([]int64, 0, len(articles))
	for _, a := range articles {
		articleIDs = append(articleIDs, a.ID)
		if a.CategoryID != 0 {
			categoryIDs = append(categoryIDs, a.CategoryID)
		}
	}

	categoryNames, err := s.categoryNameMap(categoryIDs)
	if err != nil {
		return err
	}
	articleTags, err := s.articleTagsMap(articleIDs)
	if err != nil {
		return err
	}

	// 填充关联字段。
	for i := range articles {
		articles[i].Date = dateOnly(articles[i].Date)
		articles[i].CategoryName = categoryNames[articles[i].CategoryID]
		articles[i].Tags = articleTags[articles[i].ID]
		if articles[i].Tags == nil {
			articles[i].Tags = []string{}
		}
		articles[i].Summary = s.summary(articles[i].Path)
	}
	return nil
}

// summaryRunes 是列表预览提取的最大字符数（足够填满约四行）。
const summaryRunes = 200

// summary 读取 html 文件并提取正文纯文本预览，读取失败时返回空串。
func (s *ArticleService) summary(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return common.HTMLToPreview(content, summaryRunes)
}

// categoryNameMap 批量查询分类，返回 category_id -> name 的映射。
func (s *ArticleService) categoryNameMap(categoryIDs []int64) (map[int64]string, error) {
	categories, err := model.ListCategoriesByIDs(categoryIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]string, len(categories))
	for _, c := range categories {
		result[c.ID] = c.Name
	}
	return result, nil
}

// articleTagsMap 查询文章标签关联，返回 article_id -> []tag_name 的映射。
func (s *ArticleService) articleTagsMap(articleIDs []int64) (map[int64][]string, error) {
	result := make(map[int64][]string)

	relations, err := model.ListArticleTagsByArticleIDs(articleIDs)
	if err != nil {
		return nil, err
	}
	if len(relations) == 0 {
		return result, nil
	}

	// 收集用到的标签 id，批量查出标签名。
	tagIDs := make([]int64, 0, len(relations))
	for _, r := range relations {
		tagIDs = append(tagIDs, r.TagID)
	}
	tags, err := model.ListTagsByIDs(tagIDs)
	if err != nil {
		return nil, err
	}
	tagNames := make(map[int64]string, len(tags))
	for _, t := range tags {
		tagNames[t.ID] = t.Name
	}

	// 按文章聚合标签名。
	for _, r := range relations {
		name, ok := tagNames[r.TagID]
		if !ok {
			continue
		}
		result[r.ArticleID] = append(result[r.ArticleID], name)
	}
	return result, nil
}
