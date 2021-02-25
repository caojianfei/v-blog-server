package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"v-blog/config"
)

type File struct {
	gorm.Model
	Name string
	Date string
	FullName string
	Md5 string
	MimeType string
	Ext string
	Size int64
}

// 获取文件公开访问地址
func (f *File) Url() (string, error) {
	conf, err := config.Get()
	if err != nil {
		return "", err
	}

	if f.ID == 0 {
		return "", errors.New("文件不存在")
	}

	if f.IsImage() {
		return fmt.Sprintf("%s/images/%s/%s", conf.App.Host, f.Date, f.Name), nil
	}

	return "", errors.New("非法访问")
}

// 判断文件是否是图片
func (f *File) IsImage() bool {
	types := []string{
		"image/jpeg", "image/gif", "image/png",
	}

	for _, t := range types {
		if t == f.MimeType {
			return true
		}
	}

	return false
}
