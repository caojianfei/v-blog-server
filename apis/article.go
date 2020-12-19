package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/fatih/set.v0"
	"strconv"
	"time"
	"v-blog/consts"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/models"
)

type ArticleController struct {
}

var Article ArticleController

// 文章列表
func (c ArticleController) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		category := c.Query("category")
		tag := c.Query("tag")
		title := c.Query("title")
		page := c.Query("page")
		pageSize := c.Query("pageSize")

		categoryId, _ := strconv.Atoi(category)
		tagId, _ := strconv.Atoi(tag)
		currentPage, _ := strconv.Atoi(page)
		count, _ := strconv.Atoi(pageSize)

		if count <= 0 {
			count = 10
		}
		if currentPage < 1 {
			currentPage = 1
		}

		var articleIds []int
		if tagId > 0 {
			databases.DB.Table("article_tags").Where("tag_id = ?", tagId).Pluck("article_id", &articleIds)
		}

		query := databases.DB.Model(&models.Article{}).Where("is_draft = ?", 0).Where("published_at <= ?", time.Now()).Order("views desc, id desc")
		if categoryId > 0 {
			query = query.Where("category_id = ?", categoryId)
		}
		if len(articleIds) > 0 {
			query = query.Where("id IN (?)", articleIds)
		}

		if title != "" {
			query = query.Where("title like ?", fmt.Sprintf("%%%s%%", title))
		}

		var articles []models.Article
		var total int
		query.Count(&total)
		query.Preload("Category").Preload("Tags").Offset((currentPage - 1) * count).Limit(count).Find(&articles)

		// 文章图片 md5 切片
		articleImageMd5 := make([]string, 0)
		articleIdSet := set.New(set.ThreadSafe)
		articleImageMd5Set := set.New(set.ThreadSafe)

		for _, a := range articles {
			articleIdSet.Add(a.ID)
			if a.HeadImage != "" {
				articleImageMd5Set.Add(a.HeadImage)
			}
		}

		// 文章评论数量
		var commentCounts []struct {
			ArticleId    uint
			CommentCount uint
		}
		commentCountsMap := make(map[uint]uint)
		databases.DB.
			Table("comments").
			Select("article_id, count(*) as comment_count").
			Where("article_id in (?)", articleIdSet.List()).
			Where("state = ?", 1).
			Group("article_id").
			Find(&commentCounts)
		for _, item := range commentCounts {
			commentCountsMap[item.ArticleId] = item.CommentCount
		}
		// 文章图片
		for _, item := range articleImageMd5Set.List() {
			switch val := item.(type) {
			case string:
				articleImageMd5 = append(articleImageMd5, val)
			}
		}
		imgMap, _ := helpers.BatchGetImageUrlsByMd5(articleImageMd5)

		list := make([]gin.H, len(articles))
		for k, article := range articles {
			headImageUrl := ""
			if article.HeadImage != "" {
				if url, ok := imgMap[article.HeadImage]; ok {
					headImageUrl = url
				}
			}

			articleTags := make([]gin.H, 0)
			for _, tag := range article.Tags {
				articleTags = append(articleTags, gin.H{
					"value": tag.ID,
					"label": tag.Name,
				})
			}

			var commentNum uint
			commentNum, _ = commentCountsMap[article.ID]

			list[k] = gin.H{
				"id":           article.ID,
				"title":        article.Title,
				"headImageUrl": headImageUrl,
				"intro":        article.Intro,
				"views":        article.Views,
				"commentCount": commentNum,
				"category": gin.H{
					"label": article.Category.Name,
					"value": article.Category.ID,
				},
				"publishedAt": article.PublishedAt.Format(consts.DefaultTimeFormat),
				"tags":        articleTags,
			}
		}

		helpers.ResponseOk(c, "success", &gin.H{
			"list":        list,
			"total":       total,
			"currentPage": currentPage,
		})
	}
}

// 文章详情
func (c ArticleController) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		preview := c.Query("preview")
		isPreview, _ := strconv.Atoi(preview)

		id, err := helpers.GetIdFromParam(c)
		if err != nil {
			helpers.ResponseError(c, helpers.RequestParamError, "参数错误")
		}

		var article models.Article
		query := databases.DB.Preload("Category").Preload("Tags")
		if isPreview != 1 {
			query = query.Where("is_draft = ?", 0).Where("published_at <= ?", time.Now())
		}
		err = query.First(&article, id).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				helpers.ResponseError(c, helpers.RecordNotFound, "文章不存在")
				return
			} else {
				helpers.ResponseError(c, helpers.DatabaseUnknownErr, "查询失败")
				return
			}
		}

		// 增加浏览量
		if isPreview != 1 {
			article.IncreaseViewCount()
		}

		result := gin.H{
			"id":           article.ID,
			"title":        article.Title,
			"headImageUrl": helpers.SingleGetImageUrlByMd5(article.HeadImage),
			"content":      article.Content,
			"intro":        article.Intro,
			"category": gin.H{
				"value": article.Category.ID,
				"label": article.Category.Name,
			},
			"views":        article.Views,
			"commentCount": article.CommentCount,
			"publishedAt":  article.PublishedAt.Format(consts.DefaultTimeFormat),
		}

		articleTags := make([]gin.H, 0, 5)
		for _, tag := range article.Tags {
			articleTags = append(articleTags, gin.H{
				"value": tag.ID,
				"label": tag.Name,
			})
		}

		result["tags"] = articleTags

		helpers.ResponseOk(c, "success", &result)
	}
}
