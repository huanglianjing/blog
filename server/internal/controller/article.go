package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/huanglianjing/blog/server/internal/common"
	"github.com/huanglianjing/blog/server/internal/service"
)

// ArticleController 处理文章相关的 HTTP 请求。
type ArticleController struct {
	svc *service.ArticleService
}

// NewArticleController 构造 ArticleController。
func NewArticleController() *ArticleController {
	return &ArticleController{svc: service.NewArticleService()}
}

// List 处理 GET /article/list，按页返回文章列表。
// 查询参数 page 从 0 开始，缺省或非法时按 0 处理。
func (c *ArticleController) List(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	if err != nil || page < 0 {
		page = 0
	}

	result, err := c.svc.List(page)
	if err != nil {
		common.Fail(ctx, 1, err.Error())
		return
	}
	common.OK(ctx, result)
}

// Detail 处理 GET /article/detail，按标题返回文章详情。
// 查询参数 title 为文章标题。
func (c *ArticleController) Detail(ctx *gin.Context) {
	title := ctx.Query("title")
	if title == "" {
		common.Fail(ctx, 1, "缺少参数 title")
		return
	}

	result, err := c.svc.Detail(title)
	if err != nil {
		common.Fail(ctx, 1, err.Error())
		return
	}
	if result == nil {
		common.Fail(ctx, 2, "文章不存在")
		return
	}
	common.OK(ctx, result)
}
