package apis

import (
	"github.com/gin-gonic/gin"
	"v-blog/consts"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/models"
)

type CommentController struct {
}

var Comment CommentController

func (c CommentController) List() gin.HandlerFunc  {
	return func(c *gin.Context) {
		articleId := c.Query("articleId")
		if articleId == "" {
			 helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			 return
		}

		var comments []models.Comment
		databases.DB.Where("article_id = ?", articleId).Where("state = ?", 1).Find(&comments)

		result := make([]gin.H, 0, len(comments))
		for _, comment := range comments{
			result = append(result, gin.H{
				"id": comment.ID,
				"nickname": comment.Nickname,
				"content": comment.Content,
				"createdAt": comment.CreatedAt.Format(consts.DefaultTimeFormat),
			})
		}

		helpers.ResponseOk(c, "success", &gin.H{
			"list": result,
		})
	}
}

type CommentEdit struct {
	Nickname string  `form:"nickname" binding:"required"`
	Email string	`form:"email"`
	Content string `form:"content" binding:"required"`
	ArticleId uint `form:"articleId" binding:"required"`
}

func (c CommentController) Create() gin.HandlerFunc  {
	return func(c *gin.Context) {
		var form CommentEdit
		if err := c.ShouldBind(&form); err != nil {
			helpers.ResponseValidateError(c, err)
			return
		}

		comment := models.Comment{
			Nickname: form.Nickname,
			Email: form.Email,
			Content: form.Content,
			ArticleId: form.ArticleId,
		}

		databases.DB.Create(&comment)

		if comment.ID > 0 {
			helpers.ResponseOkWithoutData(c, "success")
			return
		}

		helpers.ResponseError(c, helpers.RecordCreatedFail, "评论失败")
		return
	}
}

