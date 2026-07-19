package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/huanglianjing/blog/server/internal/common"
	"github.com/huanglianjing/blog/server/internal/service"
)

// TagController 处理标签相关的 HTTP 请求。
type TagController struct {
	svc *service.TagService
}

// NewTagController 构造 TagController。
func NewTagController() *TagController {
	return &TagController{svc: service.NewTagService()}
}

// Overview 处理 GET /tag/overview，返回各标签及其文章数（按文章数降序）。
func (c *TagController) Overview(ctx *gin.Context) {
	result, err := c.svc.Overview()
	if err != nil {
		common.Fail(ctx, 1, err.Error())
		return
	}
	common.OK(ctx, result)
}

// List 处理 GET /tag/list，按标签名分页返回文章列表。
// 查询参数 name 为标签名称，page 从 0 开始，缺省或非法时按 0 处理。
func (c *TagController) List(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		common.Fail(ctx, 1, "缺少参数 name")
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	if err != nil || page < 0 {
		page = 0
	}

	result, err := c.svc.List(name, page)
	if err != nil {
		common.Fail(ctx, 1, err.Error())
		return
	}
	if result == nil {
		common.Fail(ctx, 2, "标签不存在")
		return
	}
	common.OK(ctx, result)
}
