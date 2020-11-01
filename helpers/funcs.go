package helpers

import (
	"errors"
	"github.com/jinzhu/gorm"
	"os"
	"v-blog/databases"
	"v-blog/models"
)

// 判断文件是否存在
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

func getImageUrlByMd5(md5 string) (string, error) {
	file := models.File{}
	err := databases.DB.Where("md5", md5).First(&file).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return "", errors.New("图片不存在")
		} else {
			return "", errors.New("获取图片失败")
		}
	}

	return file.Url()
}
