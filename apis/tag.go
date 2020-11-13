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

		result := make([]gin.H, 0, len(tags))
		for _, tag := range tags {
			result = append(result, gin.H{
				"name":         tag.Name,
				"id":           tag.ID,
				"articleCount": tag.ArticleCount,
			})
		}

		helpers.ResponseOk(c, "success", &gin.H{
			"list": result,
		})
		return
	}
}
