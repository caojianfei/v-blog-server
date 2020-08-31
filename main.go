package main

import (
	"fmt"
	"github.com/spf13/viper"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/routers"
)

func init()  {
	// 初始化配置
	initConfig()
	// 初始化数据库
	databases.InitDatabase()
}

func main() {
	helpers.InitValidator()
	 _ = routers.Router.Run(":8888")
}


func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.WatchConfig()
}
