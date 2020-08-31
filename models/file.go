package models

import "github.com/jinzhu/gorm"

type File struct {
	gorm.Model
	Name string
	FullName string
	Md5 string
	MimeType string
	Ext string
	Size int64
}
