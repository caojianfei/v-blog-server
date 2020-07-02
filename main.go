package main

import (
	"fmt"
	_ "v-blog/routers"
)

func main() {

	fmt.Println("hello world!")

	//helpers.InitValidator()

	// _ = routers.Router.Run(":8888")

}
