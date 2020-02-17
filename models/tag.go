package models

import (
	"github.com/jinzhu/gorm"
	"v-blog/databases"
)

type Tag struct {
	gorm.Model
	Name string
	Description string
	ArticleCount uint
}

func (tag *Tag) IncreaseArticleCount() {
	databases.DB.Model(tag).UpdateColumn("article_count", gorm.Expr("article_count + ?", 1))
}

func (tag *Tag) DecreaseArticleCount() {
	if tag.ArticleCount == 0 {
		return
	}
	databases.DB.Model(tag).UpdateColumn("article_count", gorm.Expr("article_count - ?", 1))
}