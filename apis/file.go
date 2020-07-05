package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"v-blog/helpers"
)

const (
	rootPath  = "./files"
	imagePath = rootPath + "/images"
)

type FileController struct {
}

var File FileController

func (c FileController) UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			fmt.Printf("image upload error: %s\n", err)
			return
		}

		savePath := imagePath
		pathExists, err := PathExists(savePath)
		if err != nil {
			helpers.ResponseError(c, helpers.PathBaseError, "上传失败")
			return
		}
		fmt.Println("savePath", savePath)

		if pathExists == false {
			if err = os.Mkdir(savePath, 0755); err != nil {
				helpers.ResponseError(c, helpers.PathCreateFail, "上传失败")
				return
			}
		}

		err = c.SaveUploadedFile(file, savePath + "/" + file.Filename)
		if err != nil {
			fmt.Println("image save error", err)
		}

		log.Println(file.Filename)
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
