package main

import (
	"v-blog/helpers"
	"v-blog/routers"
	_ "v-blog/routers"
	"github.com/spf13/viper"
)

func main() {
	helpers.InitValidator()
	 _ = routers.Router.Run(":8888")
}


func initConfig() {
	viper.SetConfigName("config")
}
