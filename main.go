package main

import (
	"v-blog/helpers"
	"v-blog/routers"
	_ "v-blog/routers"
)

func main() {
	helpers.InitValidator()
	_ = routers.Router.Run(":8888")

}
