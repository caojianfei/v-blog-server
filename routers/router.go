package routers

import (
	"github.com/gin-gonic/gin"
	"v-blog/apis"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
	Router.Use(Cors())
	registerAdminRoute()
}

func registerAdminRoute() {
	admin := Router.Group("/admin")
	{
		admin.POST("/login", apis.Login)
		needLogin := admin.Use(CheckLogin())
		{
			// 文章接口
			needLogin.GET("/articles", apis.Article.List())
			needLogin.POST("/article", apis.Article.Create())
			needLogin.GET("/article/:id", apis.Article.Show())
			needLogin.POST("/article/:id", apis.Article.Edit())
			needLogin.DELETE("article/:id", apis.Article.Delete())

			// 分类接口
			needLogin.POST("/category", apis.Category.Create())
			needLogin.GET("/categories", apis.Category.List())
			needLogin.DELETE("/category/:id", apis.Category.Delete())
			needLogin.POST("/category/:id", apis.Category.Edit())
			needLogin.GET("/category/:id", apis.Category.Show())

			// 标签接口
			needLogin.POST("/tag", apis.Tag.Create())
			needLogin.POST("/tag/:id", apis.Tag.Edit())
			needLogin.GET("/tag/:id", apis.Tag.Show())
			needLogin.GET("/tags", apis.Tag.List())
			needLogin.DELETE("/tag/:id", apis.Tag.Delete())
		}
	}
}

