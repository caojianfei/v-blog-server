package main

import (
	"fmt"
	"v-blog/config"
)

func main() {
	//config.InitConfig(&config.Param{})
	//config.InitConfig(&config.Param{})
	//c, err := config.Get()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//
	//fmt.Println(c.Db.Host)

	for i := 0; i < 5; i++ {
		go func(i int) {
			fmt.Println("线程：", i)
			config.InitConfig(&config.Param{})
			fmt.Println("线程：", i, "执行完成")
		}(i)
	}

	for {
		// 发送大量进口了
	}

}