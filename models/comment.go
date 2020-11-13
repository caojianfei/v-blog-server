package models

import "github.com/jinzhu/gorm"

type Comment struct {
	gorm.Model
	ArticleId uint
	Article Article
	Nickname string
	Email string
	Content string
	State int
}
