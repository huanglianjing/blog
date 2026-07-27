package service

import (
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/huanglianjing/blog/server/internal/common"
	"github.com/huanglianjing/blog/server/internal/model"
)

// SearchService 承载搜索相关的业务逻辑。
type SearchService struct {
	articleSvc *ArticleService
}

// NewSearchService 构造 SearchService。
func NewSearchService() *SearchService {
	return &SearchService{articleSvc: NewArticleService()}
}

// 搜索结果的四个大类，对应 /search/list 接口的 type 参数。
const (
	SearchTypeCategory = "category"
	SearchTypeTag      = "tag"
	SearchTypeTitle    = "title"
	SearchTypeContent  = "content"
)

// OverviewSize 是概览接口每个大类返回的最大记录数。
const OverviewSize = 5

// SearchGroup 是概览中一个大类的结果：List 最多 OverviewSize 条，
// Total 是该类命中总数，前端据此判断是否展示「更多」入口。
type SearchGroup[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}

// NameCount 是分类 / 标签的命中结果，两者结构一致故统一返回。
type NameCount struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// SearchOverviewResult 是搜索概览接口的返回结构，四个大类按展示顺序排列。
type SearchOverviewResult struct {
	Keyword    string                     `json:"keyword"`
	Categories SearchGroup[NameCount]     `json:"categories"`
	Tags       SearchGroup[NameCount]     `json:"tags"`
	Titles     SearchGroup[model.Article] `json:"titles"`
	Contents   SearchGroup[model.Article] `json:"contents"`
}

// SearchNameListResult 是分类 / 标签类「更多」页的返回结构，一次返回全部结果。
type SearchNameListResult struct {
	List []NameCount `json:"list"`
}

// Overview 返回四个大类各自的前 OverviewSize 条命中结果与命中总数。
func (s *SearchService) Overview(keyword string) (*SearchOverviewResult, error) {
	categories, err := s.searchCategories(keyword)
	if err != nil {
		return nil, err
	}
	tags, err := s.searchTags(keyword)
	if err != nil {
		return nil, err
	}

	titleTotal, err := model.CountArticlesByTitle(keyword)
	if err != nil {
		return nil, err
	}
	titles, err := s.searchArticles(keyword, SearchTypeTitle, 0, OverviewSize)
	if err != nil {
		return nil, err
	}

	contentTotal, err := model.CountArticlesByContent(keyword)
	if err != nil {
		return nil, err
	}
	contents, err := s.searchArticles(keyword, SearchTypeContent, 0, OverviewSize)
	if err != nil {
		return nil, err
	}

	return &SearchOverviewResult{
		Keyword: keyword,
		Categories: SearchGroup[NameCount]{
			List: truncate(categories, OverviewSize), Total: int64(len(categories)),
		},
		Tags: SearchGroup[NameCount]{
			List: truncate(tags, OverviewSize), Total: int64(len(tags)),
		},
		Titles:   SearchGroup[model.Article]{List: titles, Total: titleTotal},
		Contents: SearchGroup[model.Article]{List: contents, Total: contentTotal},
	}, nil
}

// NameList 返回分类或标签类的全部命中结果，不分页。
// typ 只接受 SearchTypeCategory / SearchTypeTag。
func (s *SearchService) NameList(keyword, typ string) (*SearchNameListResult, error) {
	var list []NameCount
	var err error
	if typ == SearchTypeCategory {
		list, err = s.searchCategories(keyword)
	} else {
		list, err = s.searchTags(keyword)
	}
	if err != nil {
		return nil, err
	}
	return &SearchNameListResult{List: list}, nil
}

// ArticleList 按页返回标题或正文类的命中文章，每页 PageSize 条。
// typ 只接受 SearchTypeTitle / SearchTypeContent。
func (s *SearchService) ArticleList(keyword, typ string, page int) (*ArticleListResult, error) {
	if page < 0 {
		page = 0
	}

	var total int64
	var err error
	if typ == SearchTypeTitle {
		total, err = model.CountArticlesByTitle(keyword)
	} else {
		total, err = model.CountArticlesByContent(keyword)
	}
	if err != nil {
		return nil, err
	}
	totalPages := int((total + int64(PageSize) - 1) / int64(PageSize))

	articles, err := s.searchArticles(keyword, typ, page*PageSize, PageSize)
	if err != nil {
		return nil, err
	}
	return &ArticleListResult{List: articles, TotalPages: totalPages}, nil
}

// searchCategories 按分类名子串匹配，沿用分类概览的文章数降序。
func (s *SearchService) searchCategories(keyword string) ([]NameCount, error) {
	list, err := model.SearchCategoriesWithCount(keyword)
	if err != nil {
		return nil, err
	}
	result := make([]NameCount, 0, len(list))
	for _, c := range list {
		result = append(result, NameCount{Name: c.Name, Count: c.Count})
	}
	return result, nil
}

// searchTags 按标签名子串匹配，沿用标签概览的名称排序规则。
func (s *SearchService) searchTags(keyword string) ([]NameCount, error) {
	list, err := model.SearchTagsWithCount(keyword)
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

	result := make([]NameCount, 0, len(list))
	for _, t := range list {
		result = append(result, NameCount{Name: t.Name, Count: t.Count})
	}
	return result, nil
}

// searchArticles 按标题或正文查询文章，填充关联字段并截取命中片段。
func (s *SearchService) searchArticles(keyword, typ string, offset, limit int) ([]model.Article, error) {
	var articles []model.Article
	var err error
	if typ == SearchTypeTitle {
		articles, err = model.SearchArticlesByTitle(keyword, offset, limit)
	} else {
		articles, err = model.SearchArticlesByContent(keyword, offset, limit)
	}
	if err != nil {
		return nil, err
	}
	if err := s.articleSvc.enrichArticles(articles); err != nil {
		return nil, err
	}

	// 正文也命中时给出片段；仅标题命中的文章片段为空，前端回落到摘要。
	for i := range articles {
		articles[i].Snippet = snippet(articles[i].Content, keyword)
	}
	return articles, nil
}

// 命中片段的截取长度：命中位置前保留 snippetBefore 个字符，整段共 snippetTotal 个字符。
const (
	snippetBefore = 30
	snippetTotal  = 120
)

// snippet 从正文中截取首次命中关键词处的上下文片段，截断处补省略号。
// 未命中或正文为空时返回空串。与 SQLite 的 LIKE 一致地忽略大小写。
func snippet(content, keyword string) string {
	if content == "" || keyword == "" {
		return ""
	}

	// strings.ToLower 逐个 rune 转换，字符数不变，
	// 故小写副本中的 rune 下标可直接用于原文。
	lower := strings.ToLower(content)
	byteIndex := strings.Index(lower, strings.ToLower(keyword))
	if byteIndex < 0 {
		return ""
	}
	hit := utf8.RuneCountInString(lower[:byteIndex])

	runes := []rune(content)
	start := max(hit-snippetBefore, 0)
	end := min(start+snippetTotal, len(runes))

	var b strings.Builder
	if start > 0 {
		b.WriteString("…")
	}
	b.WriteString(string(runes[start:end]))
	if end < len(runes) {
		b.WriteString("…")
	}
	return b.String()
}

// truncate 返回切片的前 n 个元素，不足 n 个时原样返回。
func truncate[T any](list []T, n int) []T {
	if len(list) > n {
		return list[:n]
	}
	return list
}
