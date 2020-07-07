package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {
	s := gin.Default()
	s.Static("/assets", "./images/upload")
	s.POST("/upload", handler())
	s.Run(":9999")
}


func handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := "./images/upload/tmp"
		form, _ := c.MultipartForm()
		files, ok := form.File["files"]
		if exist, _ := PathExists(path); !exist {
			fmt.Printf("path: {%s} not exist", path)
			err := os.MkdirAll(path, 0755)
			if err != nil {
				log.Fatalf("create path err: %s\n", err)
			}
		}
		if ok {
			fmt.Println("ok")
			for _, file := range files {
				err := c.SaveUploadedFile(file, path + "/" + file.Filename)
				//fmt.Println("err", err)
				if err != nil {
					log.Fatalf("SaveUploadedFile err: %s\n", err)
				}
			}
		}
		//fmt.Println(files, ok)

		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "success"})
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
