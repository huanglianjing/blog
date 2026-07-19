package router

import (
	"github.com/gin-gonic/gin"

	"github.com/huanglianjing/blog/server/internal/controller"
)

// New 构造并注册所有路由。
func New() *gin.Engine {
	engine := gin.Default()

	articleCtrl := controller.NewArticleController()
	categoryCtrl := controller.NewCategoryController()
	tagCtrl := controller.NewTagController()

	article := engine.Group("/article")
	{
		article.GET("/list", articleCtrl.List)
		article.GET("/detail", articleCtrl.Detail)
	}

	category := engine.Group("/category")
	{
		category.GET("/overview", categoryCtrl.Overview)
		category.GET("/list", categoryCtrl.List)
	}

	tag := engine.Group("/tag")
	{
		tag.GET("/overview", tagCtrl.Overview)
		tag.GET("/list", tagCtrl.List)
	}

	return engine
}
