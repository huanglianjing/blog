package controller

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/huanglianjing/blog/server/internal/common"
	"github.com/huanglianjing/blog/server/internal/service"
)

// SearchController 处理搜索相关的 HTTP 请求。
type SearchController struct {
	svc *service.SearchService
}

// NewSearchController 构造 SearchController。
func NewSearchController() *SearchController {
	return &SearchController{svc: service.NewSearchService()}
}

// maxKeywordRunes 是关键词的最大字符数，超出部分截断。
const maxKeywordRunes = 100

// Overview 处理 GET /search/overview，返回分类 / 标签 / 标题 / 正文四类的命中概览。
// 查询参数 keyword 为搜索关键词。
func (c *SearchController) Overview(ctx *gin.Context) {
	keyword, ok := queryKeyword(ctx)
	if !ok {
		return
	}

	result, err := c.svc.Overview(keyword)
	if err != nil {
		common.Fail(ctx, 1, err.Error())
		return
	}
	common.OK(ctx, result)
}

// List 处理 GET /search/list，返回某一大类的完整命中结果。
// 查询参数 keyword 为关键词，type 为大类；分类 / 标签一次返回全部，
// 标题 / 正文按 page 分页（从 0 开始）。
func (c *SearchController) List(ctx *gin.Context) {
	keyword, ok := queryKeyword(ctx)
	if !ok {
		return
	}

	typ := ctx.Query("type")
	switch typ {
	case service.SearchTypeCategory, service.SearchTypeTag:
		result, err := c.svc.NameList(keyword, typ)
		if err != nil {
			common.Fail(ctx, 1, err.Error())
			return
		}
		common.OK(ctx, result)
	case service.SearchTypeTitle, service.SearchTypeContent:
		page, err := strconv.Atoi(ctx.DefaultQuery("page", "0"))
		if err != nil || page < 0 {
			page = 0
		}
		result, err := c.svc.ArticleList(keyword, typ, page)
		if err != nil {
			common.Fail(ctx, 1, err.Error())
			return
		}
		common.OK(ctx, result)
	default:
		common.Fail(ctx, 1, "参数 type 非法")
	}
}

// queryKeyword 取出并清洗关键词，为空时写入错误响应并返回 false。
func queryKeyword(ctx *gin.Context) (string, bool) {
	keyword := common.CleanName(ctx.Query("keyword"))
	if keyword == "" {
		common.Fail(ctx, 1, "缺少参数 keyword")
		return "", false
	}
	if runes := []rune(keyword); len(runes) > maxKeywordRunes {
		keyword = strings.TrimSpace(string(runes[:maxKeywordRunes]))
	}
	return keyword, true
}
