package admin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"path"
	"time"
	"v-blog/config"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/models"
)

type FileController struct {
}

var File FileController

func (c FileController) UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		imageTypes := map[string]string{
			"image/jpeg": "image/jpeg",
			"image/gif": "image/gif",
			"image/png": "image/png",
		}
		form, err := c.MultipartForm()
		if err != nil {
			helpers.ResponseError(c, helpers.UploadParamInvalid, "上传失败")
			return
		}
		images := form.File["images"]
		if len(images) == 0 {
			helpers.ResponseError(c, helpers.UploadFileEmpty, "上传文件为空")
			return
		}

		result := make([]gin.H, len(images))
		for index, image := range images {
			contentType := image.Header.Get("Content-Type")
			if _, ok := imageTypes[contentType]; !ok {
				result[index] = gin.H{}
				continue
			}

			conf, _ := config.Get()
			uploadedFile, err := uploadFile(c, image, conf.UploadDir.Images)
			if err != nil {
				result[index] = gin.H{}
				return
			}

			url, err := uploadedFile.Url()
			if err != nil {
				result[index] = gin.H{}
			}

			result[index] = gin.H{
				"url": url,
				"name": uploadedFile.Name,
				"id": uploadedFile.ID,
				"md5": uploadedFile.Md5,
			}
		}

		helpers.ResponseOk(c, "success", &gin.H{"list": result})
	}
}


func uploadFile(c *gin.Context ,file *multipart.FileHeader, basePath string) (models.File, error) {
	uploadedFile := models.File{}
	f, e := file.Open()
	if e != nil {
		return uploadedFile, e
	}
	length := file.Size
	content := make([]byte, length)
	_, e = f.Read(content)
	if e != nil {
		return uploadedFile, e
	}

	h := md5.New()
	h.Write(content)
	md5Str := hex.EncodeToString(h.Sum(nil))

	databases.DB.Where("md5 = ?", md5Str).First(&uploadedFile)
	if uploadedFile.ID > 0 {
		return uploadedFile, nil
	}

	date := fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	basePath = fmt.Sprintf("%s/%s", basePath, date)
	pathExist, e := helpers.PathExists(basePath)
	if e != nil {
		return uploadedFile, e
	}
	if !pathExist {
		if e = os.MkdirAll(basePath, 0755); e != nil {
			return uploadedFile, e
		}
	}

	fileExt := path.Ext(file.Filename)
	filename := md5Str + fileExt
	savePath := basePath + "/" + filename

	e = c.SaveUploadedFile(file, savePath)
	if e != nil {
		return uploadedFile, e
	}

	uploadedFile.Size = file.Size
	uploadedFile.Name = filename
	uploadedFile.Date = date
	uploadedFile.Ext = fileExt
	uploadedFile.Md5 = md5Str
	uploadedFile.FullName = savePath
	uploadedFile.MimeType = file.Header.Get("Content-Type")

	e = databases.DB.Create(&uploadedFile).Error
	if e != nil {
		return uploadedFile, e
	}

	return uploadedFile, nil
}
