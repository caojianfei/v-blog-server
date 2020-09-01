package main

import (
	"v-blog/config"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/routers"
)


func init()  {
	// 初始化配置
	config.InitConfig(&config.Param{})
	// 初始化数据库
	databases.InitDatabase()
}

func main() {
	helpers.InitValidator()
	_ = routers.Router.Run(":8888")
}
