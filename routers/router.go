package routers

import (
	"github.com/gin-gonic/gin"
	"log"
	"v-blog/apis"
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
			needLogin.GET("/categories/:title", apis.Category.QueryByName())

			// 标签接口
			needLogin.POST("/tag", apis.Tag.Create())
			needLogin.POST("/tag/:id", apis.Tag.Edit())
			needLogin.GET("/tag/:id", apis.Tag.Show())
			needLogin.GET("/tags", apis.Tag.List())
			needLogin.DELETE("/tag/:id", apis.Tag.Delete())
		}
	}
}

func registerFileRoute() {
	conf, _ := config.Get()
	Router.POST("/files/image", apis.File.UploadImage())
	Router.Static("/images", conf.UploadDir.Images)
}
