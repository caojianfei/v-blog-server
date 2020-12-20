package apis

import (
	"github.com/gin-gonic/gin"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/models"
)

type TagController struct {
}

var Tag TagController

func (c TagController) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tags []models.Tag
		databases.DB.Find(&tags)

		articleCountMap := make(map[uint]uint)
		{
			counts := make([]struct {
				TagId        uint
				ArticleCount uint
			}, 0)
			databases.DB.Table("article_tags").Select("tag_id, count(article_id) as article_count").Group("tag_id").Find(&counts)
			for _, item := range counts {
				articleCountMap[item.TagId] = item.ArticleCount
			}
		}

		result := make([]gin.H, 0, len(tags))
		for _, tag := range tags {
			articleCount, _ := articleCountMap[tag.ID]
			result = append(result, gin.H{
				"name":         tag.Name,
				"id":           tag.ID,
				"articleCount": articleCount,
			})
		}

		helpers.ResponseOk(c, "success", &gin.H{
			"list": result,
		})
		return
	}
}
