package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"v-blog/databases"
)

type Article struct {
	gorm.Model
	Title        string
	HeadImage    string
	Content      string
	Intro        string
	Views        int
	CommentCount int
	IsDraft      int
	PublishedAt  time.Time
	CategoryId   uint
	Category     Category
	Tags         []Tag `gorm:"many2many:article_tags;"`
	Comments     []Comment
}

func (article *Article) IncreaseViewCount() {
	databases.DB.Model(&Article{}).UpdateColumn("views", gorm.Expr("views + ?", 1))
}
