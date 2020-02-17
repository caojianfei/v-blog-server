package models

import "github.com/jinzhu/gorm"

type Comment struct {
	gorm.Model
	ArticleId int
	Article Article
	Nickname string
	Email string
	Content string
	State int
}
