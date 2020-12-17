package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"v-blog/consts"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/models"
)

type CommentController struct {
}

var Comment CommentController

// 获取评论列表
func (comment CommentController) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		stateStr := c.DefaultQuery("state", "")
		pageStr := c.DefaultQuery("page", "1")
		pageSizeStr := c.DefaultQuery("pageSize", "20")
		articleIdStr := c.DefaultQuery("articleId", "0")

		var page, pageSize, total, articleId, state int
		var err error
		comments := make([]models.Comment, 0)
		query := databases.DB

		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 {
			pageSize = 20
		}
		articleId, err = strconv.Atoi(articleIdStr)
		if err == nil && articleId > 0 {
			query = query.Where("article_id = ?", articleId)
		}
		state, err = strconv.Atoi(stateStr)
		if err == nil && (state == 0 || state == 1 || state == 2) {
			query = query.Where("state = ?", state)
		}

		query.Model(&models.Comment{}).Count(&total)
		query.Preload("Article").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at desc").Find(&comments)

		result := make([]gin.H, len(comments))
		for index, item := range comments {
			result[index] = gin.H{
				"id":           item.ID,
				"articleId":    item.ArticleId,
				"nickname":     item.Nickname,
				"email":        item.Email,
				"content":      item.Content,
				"state":        item.State,
				"articleTitle": item.Article.Title,
				"createAt":     item.CreatedAt.Format(consts.DefaultTimeFormat),
				"updatedAt":    item.UpdatedAt.Format(consts.DefaultTimeFormat),
			}
		}

		helpers.ResponseOk(c, "success", &gin.H{
			"list":     result,
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
			"isEnd":    len(result) < pageSize,
		})
		return
	}
}

// 评论审核
func (comment CommentController) Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			return
		}

		stateStr := c.DefaultQuery("state", "")
		state, err := strconv.Atoi(stateStr)
		if err != nil || (state != 0 && state != 1 && state != 2) {
			helpers.ResponseError(c, helpers.RequestParamError, "参数 state 错误")
			return
		}

		var comment models.Comment
		if err = databases.DB.First(&comment, id).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				helpers.ResponseError(c, helpers.RecordNotFound, "不存在的评论")
			} else {
				helpers.ResponseError(c, helpers.DatabaseUnknownErr, "系统出错")
			}
			return
		}

		if comment.State != state {
			comment.State = state
			if err := databases.DB.Save(&comment).Error; err != nil {
				helpers.ResponseError(c, helpers.DatabaseUnknownErr, "操作失败")
				return
			}
		}

		helpers.ResponseOk(c, "success", &gin.H{})
		return
	}
}

// 删除评论
func (comment CommentController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
			return
		}

		var comment models.Comment
		if err = databases.DB.First(&comment, id).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				helpers.ResponseError(c, helpers.RecordNotFound, "该评论不存在或已删除")
			} else {
				helpers.ResponseError(c, helpers.DatabaseUnknownErr, "系统出错")
			}
			return
		}

		if err = databases.DB.Delete(&comment).Error; err != nil {
			helpers.ResponseError(c, helpers.DatabaseUnknownErr, "系统出错")
			return
		}
		helpers.ResponseOk(c, "success", &gin.H{})
		return
	}
}

