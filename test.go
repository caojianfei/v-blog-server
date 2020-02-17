package main

import (
	"fmt"
	"time"
	"v-blog/databases"
	"v-blog/models"
)

func main() {
	//data := set.New(set.ThreadSafe)
	//data.Add(1)
	//data.Add(1)
	//fmt.Println(data.Size())
	//fmt.Println(data.List())

	//var articles []models.Article
	//databases.DB.Model(models.Article{}).Select("id,title").Where([]uint{1,2,3,4}).Find(&articles)
	//fmt.Println(articles)

	//type A struct {
	//	TagId uint
	//}
	//var a []A
	//databases.DB.Table("article_tags").Select("tag_id").Group("tag_id").Find(&a)
	//fmt.Println(a)

	//rows, err := databases.DB.Table("articles").Select("title").Rows()
	//if err == nil {
	//	for rows.Next() {
	//		data := ""
	//		if err := rows.Scan(&data); err != nil {
	//			fmt.Println("err: ", err)
	//		}
	//		fmt.Println("title: ", data)
	//
	//	}
	//}

	tag := models.Tag{}
	databases.DB.Model(&models.Tag{}).First(&tag)
	// fmt.Println(tag)

	//go func() {
	//	// fmt.Println(tag
	//	// )
	//	tag.IncreaseArticleCount()
	//	fmt.Println("haha")
	//}()

	go tag.IncreaseArticleCount()
	//tag.IncreaseArticleCount()

	time.Sleep(1e9)

	//fmt.Println("In main()")
	//go longWait()
	//go shortWait()
	//fmt.Println("About to sleep in main()")
	//// sleep works with a Duration in nanoseconds (ns) !
	//time.Sleep(10 * 1e9)
	//fmt.Println("At the end of main()")
}

func longWait() {
	fmt.Println("Beginning longWait()")
	time.Sleep(5 * 1e9) // sleep for 5 seconds
	fmt.Println("End of longWait()")
}

func shortWait() {
	fmt.Println("Beginning shortWait()")
	time.Sleep(2 * 1e9) // sleep for 2 seconds
	fmt.Println("End of shortWait()")
}