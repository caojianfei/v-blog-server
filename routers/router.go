package routers

import (
	"github.com/gin-gonic/gin"
	"log"
	"v-blog/apis"
	adminApi "v-blog/apis/admin"
	"v-blog/config"
)

var Router *gin.Engine

// _ = routers.Router.Run(":8888")

func InitRouter() {
	conf, err := config.Get()
	if err != nil {
		log.Fatalf("Read config err: %s", err)
	}

	var mode string

	switch conf.AppEnv {
	case "debug":
		mode = gin.DebugMode
	case "release":
		mode = gin.ReleaseMode
	case "test":
		mode = gin.TestMode
	default:
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	Router = gin.Default()
	Router.Use(Cors())
	registerAdminRoute()
	registerFileRoute()
	registerClientRoute()
}

func registerAdminRoute() {
	admin := Router.Group("/admin")
	{
		admin.POST("/login", adminApi.Login)
		needLogin := admin.Use(CheckLogin())
		{
			// 文章接口
			needLogin.GET("/articles", adminApi.Article.List())
			needLogin.POST("/article", adminApi.Article.Create())
			needLogin.GET("/article/:id", adminApi.Article.Show())
			needLogin.POST("/article/:id", adminApi.Article.Edit())
			needLogin.DELETE("article/:id", adminApi.Article.Delete())

			// 分类接口
			needLogin.POST("/category", adminApi.Category.Create())
			needLogin.GET("/categories", adminApi.Category.List())
			needLogin.DELETE("/category/:id", adminApi.Category.Delete())
			needLogin.POST("/category/:id", adminApi.Category.Edit())
			needLogin.GET("/category/:id", adminApi.Category.Show())
			needLogin.GET("/categories/:title", adminApi.Category.QueryByName())

			// 标签接口
			needLogin.POST("/tag", adminApi.Tag.Create())
			needLogin.POST("/tag/:id", adminApi.Tag.Edit())
			needLogin.GET("/tag/:id", adminApi.Tag.Show())
			needLogin.GET("/tags", adminApi.Tag.List())
			needLogin.DELETE("/tag/:id", adminApi.Tag.Delete())
			needLogin.GET("tags/:name", adminApi.Tag.QueryByName())
		}
	}
}

func registerFileRoute() {
	conf, _ := config.Get()
	Router.POST("/files/image", adminApi.File.UploadImage())
	Router.Static("/images", conf.UploadDir.Images)
}

func registerClientRoute() {
	Router.GET("/articles", apis.Article.List())
	Router.GET("/article/:id", apis.Article.Show())
	Router.GET("/comments", apis.Comment.List())
	Router.POST("/comment", apis.Comment.Create())
}
